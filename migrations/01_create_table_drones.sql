-- +migrate Up
CREATE TABLE IF NOT EXISTS drones (
    id VARCHAR(256) PRIMARY KEY,
    position_x DOUBLE PRECISION,
    position_y DOUBLE PRECISION,
    last_seen TIMESTAMP,
    deleted_at TIMESTAMP
);
