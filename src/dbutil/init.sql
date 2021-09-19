CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(64) NOT NULL PRIMARY KEY,
    password VARCHAR(64) NOT NULL,
    registeredOn TIMESTAMP NOT NULL,
    lastLogin TIMESTAMP,
    preferences JSON NOT NULL
);