CREATE DATABASE kong;
CREATE DATABASE users;
\c users
CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    email        VARCHAR UNIQUE NOT NULL,
    name         VARCHAR        NOT NULL,
    phone_number VARCHAR        NOT NULL,
    address      TEXT           NOT NULL,
    password     VARCHAR        NOT NULL,
    role         VARCHAR        NOT NULL CHECK (role in ('customer', 'restaurant_admin', 'shipper', 'admin'))
);
INSERT INTO users(email, name, phone_number, address, password, role)
VALUES ('admin@gmail.com','Admin user', '0123456789', '22 Ao Sen, Ha Dong, Ha Noi','$2a$14$3FBsf/pwa9FiE.8h5VkzKu372XGvWdkjtys9osw2Dk0XrdfxAUY0K','admin');

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE password_reset_tokens(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES USERS(id)
);

CREATE DATABASE restaurants;
\c restaurants
CREATE TABLE restaurants
(
    id      SERIAL PRIMARY KEY,
    manager_id INT UNIQUE NOT NULL,
    name    VARCHAR(255) NOT NULL,
    description TEXT,
    address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255),
    email   VARCHAR(255),
    status  VARCHAR NOT NULL CHECK (status IN ('open', 'closed'))
);

CREATE TABLE menu_items
(
    id            SERIAL PRIMARY KEY,
    restaurant_id INT     NOT NULL,
    name          VARCHAR NOT NULL,
    description   TEXT,
    price         DECIMAL NOT NULL,
    available     BOOLEAN NOT NULL,
    FOREIGN KEY   (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
);

CREATE DATABASE shippers;
\c shippers
CREATE TABLE shippers
(
    id SERIAL PRIMARY KEY,
    user_id INT  UNIQUE NOT NULL,
    status VARCHAR NOT NULL CHECK(status in ('available', 'busy', 'offline'))
);

CREATE DATABASE orders;
\c orders
CREATE TABLE orders
(
    id            SERIAL PRIMARY KEY,
    customer_id   INT     NOT NULL,
    restaurant_id INT     NOT NULL,
    shipper_id    INT     NOT NULL,
    items_price    DECIMAL NOT NULL,
    delivery_price DECIMAL NOT NULL,
    total_price   DECIMAL NOT NULL,
    status        VARCHAR NOT NULL CHECK (status IN ('created', 'cancelled', 'delivering', 'completed') )
);
CREATE TABLE order_items(
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL
);