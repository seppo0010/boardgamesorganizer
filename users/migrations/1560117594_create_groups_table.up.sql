CREATE TABLE groups (
    id serial PRIMARY KEY,
    external_id VARCHAR (255) NOT NULL,
    source VARCHAR(40) NOT NULL,
    UNIQUE (external_id, source)
);
