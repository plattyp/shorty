-- +migrate Up
CREATE TABLE public.urls (
	"id" serial PRIMARY KEY,
	"url" varchar NOT NULL COLLATE "default",
	"slug" varchar(30) NOT NULL COLLATE "default",
	"created_at" timestamp(6) WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
	"deleted_at" timestamp(6) WITH TIME ZONE
);

CREATE UNIQUE INDEX "index_urls_unique_slug"
ON "public"."urls" ("slug");

-- +migrate Down
DROP TABLE IF EXISTS public.urls;