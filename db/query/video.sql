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