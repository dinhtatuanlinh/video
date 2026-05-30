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
	Limit    int32
	Offset   int32
	OrderBy  string
	OrderDir string
}

var allowedOrderBy = map[string]bool{
	"video_id":            true,
	"video_category_name": true,
	"created_at":          true,
}

var allowedOrderDir = map[string]bool{
	"asc":  true,
	"desc": true,
}

func (u *UseCaseVideo) GetVideos(ctx context.Context, req *GetVideosModel) ([]db.Video, error) {
	orderBy := req.OrderBy
	if !allowedOrderBy[orderBy] {
		orderBy = "created_at"
	}
	orderDir := req.OrderDir
	if !allowedOrderDir[orderDir] {
		orderDir = "desc"
	}

	request := db.ListVideosParams{
		OrderBy:     orderBy,
		OrderDir:    orderDir,
		LimitCount:  req.Limit,
		OffsetCount: req.Offset,
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
