CREATE TABLE meetings (
    id serial PRIMARY KEY,
    group_id VARCHAR (255) NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    location VARCHAR(255) NOT NULL,
    UNIQUE (group_id)
);

CREATE TABLE attendees (
    id serial PRIMARY KEY,
    group_id VARCHAR (255) NOT NULL,
    user_id VARCHAR (255) NOT NULL,
    UNIQUE (group_id, user_id)
);
