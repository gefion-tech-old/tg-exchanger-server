CREATE TABLE IF NOT EXISTS directions_ma(
    id BIGSERIAL PRIMARY KEY,
    direction_id BIGINT REFERENCES exchange_directions(id) ON DELETE CASCADE,
    ma_id BIGINT REFERENCES merchant_autopayout(id) ON DELETE CASCADE,
    service_type INT NOT NULL,
    status BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);