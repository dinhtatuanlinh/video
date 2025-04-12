CREATE TABLE "videos"
(
    "video_id" bigserial PRIMARY KEY,
    "video_category_name" VARCHAR NOT NULL,
    "code" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL UNIQUE ,
    "url" VARCHAR NOT NULL,
    "created_at"          timestamptz NOT NULL DEFAULT (now())
);