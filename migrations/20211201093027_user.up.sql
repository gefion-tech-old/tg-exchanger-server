CREATE TABLE users (
    chat_id BIGINT NOT NULL PRIMARY KEY,
    username VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);