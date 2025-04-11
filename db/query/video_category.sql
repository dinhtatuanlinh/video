-- name: CreateVideoCategory :one
INSERT INTO video_categories (
    video_category_name,
    category_parent_name,
    created_at
) VALUES (
             $1, $2, DEFAULT
         )
    RETURNING *;

-- name: GetCategory :one
SELECT *
FROM video_categories
WHERE video_category_name = $1 LIMIT 1;