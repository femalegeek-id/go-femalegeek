-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "users" (
    "id" BIGINT NOT NULL,
    "full_name" TEXT NOT NULL,
    "nick_name" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    "job" TEXT NOT NULL,
    "company" TEXT NOT NULL,
    "email" TEXT UNIQUE,
    "phone" TEXT NULL,
    "status" user_status,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMP,
    CONSTRAINT "user_id_pkey" PRIMARY KEY ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS "users";