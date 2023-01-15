-- +migrate Up
CREATE TABLE IF NOT EXISTS pilots (
    id VARCHAR(256) PRIMARY KEY,
    first_name VARCHAR(256),
    last_name VARCHAR(256),
    phone VARCHAR(256),
    email VARCHAR(256),
    registration_time TIMESTAMP,
    drone_id VARCHAR(256),
    CONSTRAINT fk_pilot_drone
        FOREIGN KEY (drone_id)
        REFERENCES drones(id)
        ON DELETE CASCADE
);
