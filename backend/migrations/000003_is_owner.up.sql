-- Владелец проекта (только один пользователь)
ALTER TABLE users ADD COLUMN is_owner BOOLEAN NOT NULL DEFAULT FALSE;

-- Гарантия: только один owner
CREATE UNIQUE INDEX idx_users_owner ON users(is_owner) WHERE is_owner = TRUE;
