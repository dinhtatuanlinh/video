// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: video.sql

package db

import (
	"context"
)

const createVideo = `-- name: CreateVideo :one
INSERT INTO videos (
    video_category_name,
    name,
    url,
    created_at
) VALUES (
             $1, $2, $3, DEFAULT
         )
    RETURNING video_id, video_category_name, name, url, created_at
`

type CreateVideoParams struct {
	VideoCategoryName string `json:"video_category_name"`
	Name              string `json:"name"`
	Url               string `json:"url"`
}

func (q *Queries) CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error) {
	row := q.db.QueryRow(ctx, createVideo, arg.VideoCategoryName, arg.Name, arg.Url)
	var i Video
	err := row.Scan(
		&i.VideoID,
		&i.VideoCategoryName,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}
