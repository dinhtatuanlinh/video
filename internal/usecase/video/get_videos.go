package video

import (
	"context"
	"errors"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"google.golang.org/grpc/codes"
	"net/http"
)

type GetVideosModel struct {
	Limit  int32
	Offset int32
}

func (u *UseCaseVideo) GetVideos(ctx context.Context, req *GetVideosModel) ([]db.Video, error) {
	request := db.ListVideosParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	videos, err := u.store.ListVideos(ctx, request)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return []db.Video{}, internalError.NewAppError("there is no video", http.StatusBadRequest, codes.NotFound, err)
		}
		return []db.Video{}, internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
	}

	return videos, nil
}
