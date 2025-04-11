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

const listVideos = `-- name: ListVideos :many
SELECT video_id, video_category_name, name, url, created_at
FROM videos
ORDER BY video_category_name LIMIT $1
OFFSET $2
`

type ListVideosParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListVideos(ctx context.Context, arg ListVideosParams) ([]Video, error) {
	rows, err := q.db.Query(ctx, listVideos, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.VideoID,
			&i.VideoCategoryName,
			&i.Name,
			&i.Url,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
