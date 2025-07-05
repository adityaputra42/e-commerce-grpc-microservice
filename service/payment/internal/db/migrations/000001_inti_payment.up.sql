-- Enable UUID extension (once per DB)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE payment_wallets (
    id SERIAL PRIMARY KEY,
    network TEXT NOT NULL UNIQUE,
    wallet_address TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    user_id UUID NOT NULL,
    network TEXT NOT NULL,                -- contoh: 'solana'
    currency TEXT NOT NULL,               -- contoh: 'USDT'
    amount NUMERIC NOT NULL,              -- jumlah yang dibayar
    wallet_address TEXT NOT NULL,         -- tujuan transfer (diambil dari payment_wallets)
    tx_hash TEXT UNIQUE,                  -- hash transaksi dari blockchain
    status TEXT NOT NULL DEFAULT 'pending', -- pending, confirmed, failed
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);


CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_network ON payments(network);
CREATE INDEX idx_payments_tx_hash ON payments(tx_hash);


