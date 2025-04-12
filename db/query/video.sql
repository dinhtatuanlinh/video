-- name: CreateVideo :one
INSERT INTO videos (
    video_category_name,
                    code,
    name,
    url,
    created_at
) VALUES (
             $1, $2, $3, $4, DEFAULT
         )
    RETURNING *;

-- name: ListVideos :many
SELECT *
FROM videos
ORDER BY video_category_name LIMIT $1
OFFSET $2;

-- name: DeleteVideoByID :exec
DELETE FROM videos
WHERE video_id = $1;

-- name: GetVideoByID :one
SELECT * FROM videos
WHERE video_id = $1;