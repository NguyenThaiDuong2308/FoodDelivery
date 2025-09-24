CREATE TABLE restaurants
(
    id      SERIAL PRIMARY KEY,
    user_id INT     NOT NULL,
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
    restaurant_id INT     NOT NULL REFERENCES restaurants (id),
    name          VARCHAR NOT NULL,
    description   TEXT,
    price         DECIMAL NOT NULL,
    available     BOOLEAN NOT NULL,
    FOREIGN KEY   (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
);