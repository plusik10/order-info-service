-- +goose Up

CREATE TABLE item (
    item_id SERIAL PRIMARY KEY,
    chrt_id INT,
    order_uid VARCHAR(255),
    track_number VARCHAR(255),
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INTEGER,
    size VARCHAR(255),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER,
    FOREIGN KEY (order_uid) REFERENCES order_doc (order_uid)
);

-- +goose Down
DROP TABLE item;