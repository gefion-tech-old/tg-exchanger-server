CREATE TABLE merchant_autopayout(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    service VARCHAR(50) NOT NULL,
    service_type INT NOT NULL,
    options TEXT,
    status BOOLEAN DEFAULT false,
    message_id BIGINT REFERENCES bot_messages(id),
    created_by VARCHAR(255) REFERENCES users(username),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);