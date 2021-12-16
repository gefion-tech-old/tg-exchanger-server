CREATE TABLE notifications(
    id BIGSERIAL NOT NULL PRIMARY KEY,    
    type INT NOT NULL,
    status INT DEFAULT 1,
    chat_id BIGINT NOT NULL,
    username VARCHAR(50) NOT NULL,
    code INT,
    user_card VARCHAR(20),
    img_path VARCHAR(255),
    ex_from VARCHAR(20),
    ex_to VARCHAR(20),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
