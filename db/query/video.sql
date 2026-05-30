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
ORDER BY
    CASE WHEN @order_by::text = 'video_id' AND @order_dir::text = 'asc' THEN video_id END ASC NULLS LAST,
    CASE WHEN @order_by::text = 'video_id' AND @order_dir::text = 'desc' THEN video_id END DESC NULLS LAST,
    CASE WHEN @order_by::text = 'video_category_name' AND @order_dir::text = 'asc' THEN video_category_name END ASC NULLS LAST,
    CASE WHEN @order_by::text = 'video_category_name' AND @order_dir::text = 'desc' THEN video_category_name END DESC NULLS LAST,
    CASE WHEN @order_by::text = 'created_at' AND @order_dir::text = 'asc' THEN created_at END ASC NULLS LAST,
    CASE WHEN @order_by::text = 'created_at' AND @order_dir::text = 'desc' THEN created_at END DESC NULLS LAST
LIMIT @limit_count
OFFSET @offset_count;

-- name: DeleteVideoByID :exec
DELETE FROM videos
WHERE video_id = $1;

-- name: GetVideoByID :one
SELECT * FROM videos
WHERE video_id = $1;