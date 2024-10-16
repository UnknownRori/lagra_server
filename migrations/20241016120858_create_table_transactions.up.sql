CREATE TABLE transactions (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    pay INT NOT NULL,
    role ENUM("CREDIT", "PAYPAL", "GOPAY") NOT NULL,

    consumer_id VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY(consumer_id) REFERENCES users(uuid)
);
