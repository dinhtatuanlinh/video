package video

import (
	"context"
	"github.com/dinhtatuanlinh/video/config"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
)

type VideoRepository interface {
	DownloadVideo(ctx context.Context, req *DownloadVideoModel) error
	CreateVideo(ctx context.Context, req *CreateVideoModel) error
	CreateVideoCategory(ctx context.Context, req *CreateVideoCategoryModel) error
	GetVideos(ctx context.Context, req *GetVideosModel) ([]db.Video, error)
	DeleteVideos(ctx context.Context, req *DeleteVideosModel) error
}
type UseCaseVideo struct {
	config config.Config
	store  db.Store
}

func NewUseCaseVideo(config config.Config, store db.Store) VideoRepository {

	return &UseCaseVideo{
		config: config,
		store:  store,
	}
}
