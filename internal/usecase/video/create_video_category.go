package video

import (
	"context"
	"errors"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"net/http"
)

type CreateVideoCategoryModel struct {
	CategoryName       string
	CategoryParentName string
}

func (u *UseCaseVideo) CreateVideoCategory(ctx context.Context, req *CreateVideoCategoryModel) error {
	if req.CategoryParentName != "" {
		_, err := u.store.GetCategory(ctx, req.CategoryParentName)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				return internalError.NewAppError("category not found", http.StatusBadRequest, codes.NotFound, err)
			}
			return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
		}
	}

	request := db.CreateVideoCategoryParams{
		CategoryParentName: req.CategoryParentName,
		VideoCategoryName:  req.CategoryName,
	}
	_, err := u.store.CreateVideoCategory(ctx, request)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			log.Error().Err(err).Msg("Error unique violation")
			return internalError.NewAppError("Error unique violation", http.StatusBadRequest, codes.InvalidArgument, err)
		}
		log.Error().Err(err).Msg("Failed to create operator transaction")
		return err
	}

	return nil
}
