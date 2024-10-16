CREATE TABLE carts (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    total INT NOT NULL,

    item_id VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY(item_id) REFERENCES items(uuid)
);
