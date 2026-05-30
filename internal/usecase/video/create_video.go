package video

import (
	"context"
	"errors"
	"fmt"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	"github.com/dinhtatuanlinh/video/internal/constant"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/dinhtatuanlinh/video/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CreateVideoModel struct {
	Videos []VideoModel
}

type VideoModel struct {
	CategoryName string
	Code         string
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
	outputDir := constant.Drive + "\\setup\\Videos"
	log.Info().Str("outputDir", outputDir).Msg("Output directory")
	for _, video := range req.Videos {
		inputPath := video.InputPath

		var fileUrl string
		var folder string
		categoryName := video.CategoryName
		for {
			category, err := u.store.GetCategory(ctx, categoryName)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					return internalError.NewAppError("category not found", http.StatusBadRequest, codes.NotFound, err)
				}
				return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
			}
			categoryName = category.CategoryParentName
			fileUrl = category.VideoCategoryName + "/" + fileUrl
			folder = category.VideoCategoryName + "/" + folder
			if category.CategoryParentName == "" {
				break
			}
		}
		folder = outputDir + "/" + folder + video.Code
		_, err := os.Stat(folder)
		if os.IsNotExist(err) {
			os.MkdirAll(folder, 0755)
		}
		fileUrl = "/videos/" + fileUrl + video.Code + "/"

		name := strings.ReplaceAll(util.RemoveSpecialCharactersButKeepSpace(video.Name), " ", "_")
		outputPath := filepath.Join(folder, name+".m3u8")

		fileName := name + ".m3u8"
		now := time.Now()
		cmd := exec.Command("ffmpeg",
			"-fflags", "+discardcorrupt", // discard corrupt frames
			"-err_detect", "ignore_err", // ignore decoding errors
			"-i", inputPath, // input file
			"-c:v", "libx264", // video codec
			"-c:a", "aac", // audio codec
			"-f", "hls", // output format
			"-hls_time", "10", // segment length in seconds
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", filepath.Join(folder, now.Format("20060102150405.000")+"_"+"%03d.ts"),
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
			Code:              video.Code,
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
