CREATE TABLE items (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL,
    img_url VARCHAR(255) NOT NULL,

    category_id VARCHAR(255) NOT NULL,
    FOREIGN KEY(category_id) REFERENCES categories(uuid) ON DELETE CASCADE
);
