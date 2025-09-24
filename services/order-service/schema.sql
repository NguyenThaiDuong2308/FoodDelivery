CREATE TABLE orders
(
    id            SERIAL PRIMARY KEY,
    customer_id   INT     NOT NULL,
    restaurant_id INT     NOT NULL,
    total_price   DECIMAL NOT NULL,
    status        VARCHAR NOT NULL CHECK (status IN
                                          ('pending', 'cancelled''accepted', 'preparing', 'delivering', 'completed') )
);

CREATE TABLE order_items(
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL NOT NULL
);