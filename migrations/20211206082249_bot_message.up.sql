CREATE TABLE bot_messages (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    connector VARCHAR(255) NOT NULL,
    message_text TEXT NOT NULL,
    created_by VARCHAR(255) REFERENCES users(username),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);