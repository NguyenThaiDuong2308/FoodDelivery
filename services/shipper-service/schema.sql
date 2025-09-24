CREATE EXTENSION IF NOT EXISTS postgis;
CREATE TABLE shippers(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status VARCHAR NOT NULL CHECK(status in ('available', 'busy', 'offline')),
    location geography(POINTm 4326) NOT NULL,
);