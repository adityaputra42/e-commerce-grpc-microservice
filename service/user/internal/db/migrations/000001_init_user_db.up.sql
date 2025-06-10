-- Aktifkan ekstensi UUID (gunakan uuid-ossp untuk uuid_generate_v4)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================
-- Tabel users
-- ==========================
CREATE TABLE users (
    username VARCHAR(50) PRIMARY KEY,  -- berasal dari auth_users.username
    full_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==========================
-- Tabel user_addresses
-- ==========================
CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL,
    label VARCHAR(100),
    recipient_name VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(20),
    address_line TEXT NOT NULL,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    is_selected BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_addresses_user FOREIGN KEY (username)
    REFERENCES users(username)
    ON DELETE CASCADE
);

-- Index untuk mempercepat pencarian berdasarkan username
CREATE INDEX idx_user_addresses_username ON user_addresses(username);

-- ==========================
-- Trigger: hanya satu alamat yang is_selected = TRUE per user
-- ==========================
CREATE OR REPLACE FUNCTION ensure_single_selected_address()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_selected THEN
        UPDATE user_addresses
        SET is_selected = FALSE
        WHERE username = NEW.username AND id <> NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_selected_address
BEFORE INSERT OR UPDATE ON user_addresses
FOR EACH ROW
EXECUTE FUNCTION ensure_single_selected_address();
