CREATE TABLE notifications(
    id BIGSERIAL NOT NULL PRIMARY KEY,    
    type INT NOT NULL,
    status INT DEFAULT 1,
    chat_id BIGINT NOT NULL,
    username VARCHAR(255) NOT NULL,
    code INT,
    user_card VARCHAR(255),
    img_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
