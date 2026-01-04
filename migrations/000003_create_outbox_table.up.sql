CREATE TABLE IF NOT EXISTS outbox (
    ID BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    EventType VARCHAR(255) NOT NULL,
    Payload TEXT NOT NULL,
    Status VARCHAR(50) NOT NULL,
    CreatedAt DATETIME NOT NULL,
    PublishedAt DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (ID),
    INDEX idx_status (Status),
    INDEX idx_created_at (CreatedAt),
    INDEX idx_status_created (Status, CreatedAt)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;