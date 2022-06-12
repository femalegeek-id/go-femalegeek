-- +migrate Up notransaction
CREATE TYPE "user_status" AS ENUM ('PENDING', 'ACTIVE', 'INACTIVE', 'BANNED');

-- +migrate Down
DROP TABLE IF EXISTS "users";