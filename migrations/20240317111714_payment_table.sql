-- +goose Up
CREATE TABLE payment (
    payment_id serial PRIMARY KEY,
    order_uid VARCHAR(255),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt INTEGER,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER,
    FOREIGN KEY (order_uid) REFERENCES order_doc (order_uid)
);
-- +goose Down
DROP TABLE payment;