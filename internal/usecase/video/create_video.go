package video

import (
	"context"
	"errors"
	"fmt"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type CreateVideoModel struct {
	Videos []VideoModel
}

type VideoModel struct {
	CategoryName string
	Name         string
	InputPath    string
}

func (u *UseCaseVideo) CreateVideo(ctx context.Context, req *CreateVideoModel) error {
	log.Info().Msg("Create Video Request")
	if len(req.Videos) < 1 {
		err := errors.New("no url provided")
		log.Error().Err(err).Msg("Failed to create operator")
		return internalError.NewAppError("Failed to DownloadVideo", http.StatusBadRequest, codes.Internal, err)
	}
	outputDir := u.config.VideoPath
	log.Info().Str("outputDir", outputDir).Msg("Output directory")
	for _, video := range req.Videos {

		_, err := os.Stat(video.InputPath)
		log.Info().Str("video", video.InputPath).Msg("Input path")
		if os.IsNotExist(err) {
			log.Error().Str("file", video.InputPath).Err(err).Msg("file not exit")
			return internalError.NewAppError("file not exit", http.StatusBadRequest, codes.Internal, err)
		}
		fileUrl := "/downloads/"
		folder := outputDir + "/"
		for {
			category, err := u.store.GetCategory(ctx, video.CategoryName)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					return internalError.NewAppError("operator_id not found", http.StatusBadRequest, codes.NotFound, err)
				}
				return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
			}
			fileUrl = fileUrl + category.VideoCategoryName + "/"
			folder = folder + category.VideoCategoryName + "/"
			if category.CategoryParentName == "" {
				break
			}
		}

		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			os.MkdirAll(folder, 0755)
		}
		outputPath := filepath.Join(folder, video.Name+".m3u8")
		fileName := video.Name + ".m3u8"
		cmd := exec.Command("ffmpeg",
			"-i", video.InputPath, // input file
			"-c:v", "libx264", // video codec
			"-c:a", "aac", // audio codec
			"-f", "hls", // output format
			"-hls_time", "10", // segment length in seconds
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", filepath.Join(folder, video.Name+"_"+"segment_%03d.ts"), // segment files
			outputPath,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("Running ffmpeg to generate .m3u8 playlist and .ts segments...")
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		request := db.CreateVideoParams{
			VideoCategoryName: video.CategoryName,
			Name:              video.Name,
			Url:               fileUrl + fileName,
		}
		_, err = u.store.CreateVideo(ctx, request)
		if err != nil {
			if db.ErrorCode(err) == db.UniqueViolation {
				log.Error().Err(err).Msg("Error unique violation")
				return internalError.NewAppError("Error unique violation", http.StatusBadRequest, codes.InvalidArgument, err)
			}
			log.Error().Err(err).Msg("Failed to create operator transaction")
			return err
		}
		fmt.Println("✅ HLS conversion complete.")
	}

	return nil
}
