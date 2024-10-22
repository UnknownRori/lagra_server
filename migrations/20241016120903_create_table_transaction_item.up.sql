CREATE TABLE transaction_item (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    total INT NOT NULL,

    item_id VARCHAR(255) NOT NULL,
    FOREIGN KEY(item_id) REFERENCES items(uuid),

    transaction_id VARCHAR(255) NOT NULL,
    FOREIGN KEY(transaction_id) REFERENCES transactions(uuid)
);
