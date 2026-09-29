-- Пользователи
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    is_superuser    BOOLEAN NOT NULL DEFAULT FALSE,
    is_admin        BOOLEAN NOT NULL DEFAULT FALSE,
    is_guide        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Маршруты
CREATE TABLE routes (
    id              BIGSERIAL PRIMARY KEY,
    owner_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    cover_photo_id  BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_routes_owner ON routes(owner_id);
CREATE INDEX idx_routes_dates ON routes(start_date, end_date);

-- Точки маршрута
CREATE TABLE places (
    id              BIGSERIAL PRIMARY KEY,
    route_id        BIGINT NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    lat             DOUBLE PRECISION NOT NULL,
    lng             DOUBLE PRECISION NOT NULL,
    description     TEXT,
    emoji           VARCHAR(10),
    order_index     INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_places_route ON places(route_id);

-- Комментарии
CREATE TABLE comments (
    id              BIGSERIAL PRIMARY KEY,
    route_id        BIGINT NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ
);

CREATE INDEX idx_comments_route ON comments(route_id);
CREATE INDEX idx_comments_user ON comments(user_id);

-- Фото
CREATE TABLE photos (
    id              BIGSERIAL PRIMARY KEY,
    route_id        BIGINT REFERENCES routes(id) ON DELETE CASCADE,
    place_id        BIGINT REFERENCES places(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_path       VARCHAR(500) NOT NULL,
    title           VARCHAR(255),
    caption         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (route_id IS NOT NULL OR place_id IS NOT NULL)
);

CREATE INDEX idx_photos_route ON photos(route_id);
CREATE INDEX idx_photos_place ON photos(place_id);
