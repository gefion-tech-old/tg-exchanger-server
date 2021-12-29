CREATE TABLE logs (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(50),
    info TEXT,
    service VARCHAR(100) NOT NULL,
    module VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT now()    
);