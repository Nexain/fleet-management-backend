-- Create the vehicle_locations table
CREATE TABLE vehicle_locations (
    vehicle_id VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp BIGINT NOT NULL,
    PRIMARY KEY (vehicle_id, timestamp)
);

-- Index for faster queries on vehicle_id and timestamp
CREATE INDEX idx_vehicle_id_timestamp ON vehicle_locations (vehicle_id, timestamp);