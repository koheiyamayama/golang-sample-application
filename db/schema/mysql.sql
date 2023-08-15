CREATE TABLE
    IF NOT EXISTS `users` (
        id CHAR(36) PRIMARY KEY,
        name VARCHAR(100) NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS `posts` (
        id CHAR(36) PRIMARY KEY,
        title VARCHAR(30) NOT NULL,
        body VARCHAR(280) NOT NULL,
        user_id CHAR(36) NOT NULL,
        INDEX idx_user_id (user_id),
        CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
    );
