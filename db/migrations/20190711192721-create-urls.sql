-- +migrate Up
CREATE TABLE public.urls (
	"id" bigserial PRIMARY KEY,
	"url" varchar NOT NULL COLLATE "default",
	"slug" varchar(30) NOT NULL COLLATE "default",
	"created_at" timestamp(6) WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
	"deleted_at" timestamp(6) WITH TIME ZONE
);

CREATE UNIQUE INDEX "index_urls_unique_slug"
	ON "public"."urls" ("slug");

CREATE INDEX "index_urls_on_created_at"
	ON "public"."urls" USING BRIN(created_at) WITH (pages_per_range = 10);

-- +migrate Down
DROP TABLE IF EXISTS public.urls;