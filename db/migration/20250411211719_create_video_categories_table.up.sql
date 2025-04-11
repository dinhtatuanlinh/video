CREATE TABLE "video_categories"
(
    "video_category_name" VARCHAR PRIMARY KEY,
    "category_parent_name" VARCHAR NOT NULL,
    "created_at"        timestamptz NOT NULL DEFAULT (now())
);