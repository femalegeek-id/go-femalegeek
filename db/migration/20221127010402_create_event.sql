-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "events" (
    "id" BIGINT NOT NULL,
    "slug" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "cover" TEXT NOT NULL,
    "organizer" TEXT NOT NULL,
    "starts_at" TIMESTAMP,
    "ends_at" TIMESTAMP,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS "events";