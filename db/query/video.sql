-- name: CreateVideo :one
INSERT INTO videos (
    video_category_name,
    name,
    url,
    created_at
) VALUES (
             $1, $2, $3, DEFAULT
         )
    RETURNING *;

-- name: ListVideos :many
SELECT *
FROM videos
ORDER BY video_category_name LIMIT $1
OFFSET $2;