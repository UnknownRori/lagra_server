CREATE TABLE users(
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM("ADMIN", "CONSUMER") NOT NULL
);