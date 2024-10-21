CREATE TABLE carts (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    total INT NOT NULL,

    user_id VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(uuid),
    item_id VARCHAR(255) NOT NULL,
    FOREIGN KEY(item_id) REFERENCES items(uuid)
);
