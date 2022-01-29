CREATE TABLE request(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    request_status INT NOT NULL,    
    exchange_from VARCHAR(20) NOT NULL,
    exchange_to VARCHAR(20) NOT NULL,
    course VARCHAR(20) NOT NULL,
    address VARCHAR(255) NOT NULL,
    expected_amount DECIMAL NOT NULL,
    transferred_amount DECIMAL DEFAULT 0,
    transaction_hash VARCHAR(255),
    created_by_username VARCHAR(255) REFERENCES users(username) NOT NULL,
    created_by_chat_id BIGINT REFERENCES users(chat_id) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);