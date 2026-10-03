CREATE SCHEMA facelab

CREATE TABLE facelab.videos (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
)

CREATE TABLE facelab.blur (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    original_video_id INTEGER NOT NULL REFERENCES facelab.videos(id)
)