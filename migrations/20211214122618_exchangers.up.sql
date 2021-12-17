CREATE TABLE exchangers (
    id BIGSERIAL NOT NULL PRIMARY KEY,    
    name VARCHAR(255) UNIQUE NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_by VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);