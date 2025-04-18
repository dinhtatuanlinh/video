// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"time"
)

type Video struct {
	VideoID           int64     `json:"video_id"`
	VideoCategoryName string    `json:"video_category_name"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	Url               string    `json:"url"`
	CreatedAt         time.Time `json:"created_at"`
}

type VideoCategory struct {
	VideoCategoryName  string    `json:"video_category_name"`
	CategoryParentName string    `json:"category_parent_name"`
	CreatedAt          time.Time `json:"created_at"`
}
