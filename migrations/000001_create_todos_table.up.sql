
CREATE TABLE IF NOT EXISTS todos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(255) NOT NULL,
    duedate TIMESTAMP NOT NULL,
    fileid VARCHAR(255) DEFAULT '',
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_createdat (createdat),
    INDEX idx_duedate (duedate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;