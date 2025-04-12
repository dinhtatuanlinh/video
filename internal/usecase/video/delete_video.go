package video

import (
	"context"
	"errors"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"net/http"
	"os"
)

type DeleteVideosModel struct {
	ID int64
}

func (u *UseCaseVideo) DeleteVideos(ctx context.Context, req *DeleteVideosModel) error {
	var videoFolder string
	video, err := u.store.GetVideoByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return internalError.NewAppError("video not found", http.StatusBadRequest, codes.NotFound, err)
		}
		return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
	}
	for {
		category, err := u.store.GetCategory(ctx, video.VideoCategoryName)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				return internalError.NewAppError("video not found", http.StatusBadRequest, codes.NotFound, err)
			}
			return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
		}
		videoFolder = category.VideoCategoryName + "/" + videoFolder
		if category.CategoryParentName == "" {
			break
		}
	}
	outputDir := u.config.VideoPath
	videoFolder = videoFolder + video.Code
	err = os.RemoveAll(outputDir + videoFolder)
	if err != nil {
		log.Error().Msgf("failed to delete video folder %s error: %s", videoFolder, err.Error())
		return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
	}
	err = u.store.DeleteVideoByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return internalError.NewAppError("there is no video", http.StatusBadRequest, codes.NotFound, err)
		}
		return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
	}

	return nil
}
