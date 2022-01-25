CREATE TABLE exchange_directions(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    exchange_from VARCHAR(20) NOT NULL,
    exchange_to VARCHAR(20) NOT NULL,
    course_correction VARCHAR(255) NOT NULL,
    directions_status BOOLEAN DEFAULT false,
    ma_service VARCHAR(50) NOT NULL,
    created_by VARCHAR(255) REFERENCES users(username),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (exchange_from, exchange_to)
);