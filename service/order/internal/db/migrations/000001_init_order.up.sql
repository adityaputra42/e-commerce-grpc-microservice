
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL,
    car_id UUID NOT NULL,
    status VARCHAR(25) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_username ON orders(username);
CREATE INDEX idx_orders_car_id ON orders(car_id);
CREATE INDEX idx_orders_status ON orders(status);
