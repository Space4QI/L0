-- Создание таблицы "orders"
CREATE TABLE IF NOT EXISTS orders
(
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(255),
    locale             VARCHAR(255),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(255),
    sm_id              INTEGER,
    date_created       TIMESTAMP WITH TIME ZONE,
    oof_shard          VARCHAR(255)
);

-- Создание таблицы "delivery"
CREATE TABLE IF NOT EXISTS delivery
(
    order_uid VARCHAR(255) PRIMARY KEY,
    name      VARCHAR(255),
    phone     VARCHAR(255),
    zip       VARCHAR(255),
    city      VARCHAR(255),
    address   VARCHAR(255),
    region    VARCHAR(255),
    email     VARCHAR(255)
);

-- Создание таблицы "payment"
CREATE TABLE IF NOT EXISTS payment
(
    order_uid     VARCHAR(255) PRIMARY KEY,
    transaction   VARCHAR(255),
    request_id    VARCHAR(255),
    currency      VARCHAR(255),
    provider      VARCHAR(255),
    amount        INTEGER,
    payment_dt    INTEGER,
    bank          VARCHAR(255),
    delivery_cost INTEGER,
    goods_total   INTEGER,
    custom_fee    INTEGER
);

-- Создание таблицы "items"
CREATE TABLE IF NOT EXISTS items
(
    order_uid   VARCHAR(255),
    chrt_id     INTEGER,
    track_number VARCHAR(255),
    price       INTEGER,
    rid         VARCHAR(255),
    name        VARCHAR(255),
    sale        INTEGER,
    size        VARCHAR(255),
    total_price INTEGER,
    nm_id       INTEGER,
    brand       VARCHAR(255),
    status      INTEGER
);


INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey,
                    sm_id, date_created, oof_shard)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '', 'test', 'meest', '9', 99, '2021-11-26T06:22:19Z',
        '1');


INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
VALUES ('b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot',
        'test@gmail.com');

INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost,
                     goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test', 'b563feb7b2b84b6test', '', 'USD', 'wbpay', 1817, 1637907727, 'alpha', 1500,
        317, 0);

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('b563feb7b2b84b6test', 9934930, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212,
        'Vivienne Sabo', 202);

