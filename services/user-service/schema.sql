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
CREATE TABLE password_reset_tokens(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES USERS(id)
);