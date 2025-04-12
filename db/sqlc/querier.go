// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
)

type Querier interface {
	CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error)
	CreateVideoCategory(ctx context.Context, arg CreateVideoCategoryParams) (VideoCategory, error)
	DeleteVideoByID(ctx context.Context, videoID int64) error
	GetCategory(ctx context.Context, videoCategoryName string) (VideoCategory, error)
	GetVideoByID(ctx context.Context, videoID int64) (Video, error)
	ListVideos(ctx context.Context, arg ListVideosParams) ([]Video, error)
}

var _ Querier = (*Queries)(nil)
