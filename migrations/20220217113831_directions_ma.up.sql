CREATE TABLE IF NOT EXISTS directions_ma(
    id BIGSERIAL ON DELETE CASCADE PRIMARY KEY,
    direction_id BIGINT REFERENCES exchange_directions(id),
    ma_id BIGINT REFERENCES merchant_autopayout(id),
    service_type INT NOT NULL,
    status BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);