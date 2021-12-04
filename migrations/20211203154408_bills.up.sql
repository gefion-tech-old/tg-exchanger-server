CREATE TABLE bills (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    chat_id BIGINT REFERENCES users(chat_id),
    bill VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now()    
);