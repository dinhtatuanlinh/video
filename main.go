package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dinhtatuanlinh/video/config"
	"github.com/dinhtatuanlinh/video/db"
	sqlc "github.com/dinhtatuanlinh/video/db/sqlc"
	deliveryapi "github.com/dinhtatuanlinh/video/internal/delivery/restful"
	"github.com/dinhtatuanlinh/video/internal/usecase"
	ucVideo "github.com/dinhtatuanlinh/video/internal/usecase/video"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"syscall"
	"time"

	"os"
	"os/signal"
)

var interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config %s", err)
	}

	if config.Environment == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	var connPool *pgxpool.Pool
	connPool, err = WaitForPostgres(config.DBSource, config.MaxRetries, 2*time.Second)
	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %s", err)
	}
	db.RunMigration(&config.MigrationUrl, &config.DBSource)

	store := sqlc.NewStore(connPool)

	waitGroup, ctx := errgroup.WithContext(ctx)

	useCaseVideo := ucVideo.NewUseCaseVideo(config, store)
	useCase := usecase.NewUseCase(useCaseVideo)
	runGinServer(ctx, waitGroup, config, useCase)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msgf("wait group finished with error: %s", err)
	}
}
func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config config.Config, usecase usecase.UseCase) {
	server, err := deliveryapi.NewServer(usecase)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	httpServer := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: server.Router,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP server at %s", config.HTTPServerAddress)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msgf("HTTP server gracefully shutting down")

		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP server gracefully")
			return err
		}

		log.Info().Msgf("HTTP server stopped")
		return nil
	})
}

func WaitForPostgres(dsn string, maxRetries int, delay time.Duration) (*pgxpool.Pool, error) {
	var dbpool *pgxpool.Pool
	var err error

	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		dbpool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			err = dbpool.Ping(ctx)
			if err == nil {
				log.Info().Msg("✅ Connected to PostgreSQL!")
				return dbpool, nil
			}
		}

		log.Warn().Err(err).Msgf("⏳ Waiting for PostgreSQL... attempt %d/%d", i, maxRetries)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("❌ Could not connect to PostgreSQL after %d attempts: %w", maxRetries, err)
}
