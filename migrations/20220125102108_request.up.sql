CREATE TABLE request(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    request_status INT NOT NULL,    
    exchange_from VARCHAR(20) NOT NULL,
    exchange_to VARCHAR(20) NOT NULL,
    course VARCHAR(20) NOT NULL,
    created_by VARCHAR(255) REFERENCES users(username),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);