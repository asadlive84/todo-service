CREATE TABLE IF NOT EXISTS files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    fileName VARCHAR(255) NOT NULL,
    originalName VARCHAR(255) NOT NULL,
    contentType VARCHAR(100) NOT NULL,
    fileSize BIGINT NOT NULL,
    filehash VARCHAR(64),
    storagepath VARCHAR(500),
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX createdat (createdat),
    INDEX idx_filename (fileName),
    INDEX idx_contentType (contentType)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;