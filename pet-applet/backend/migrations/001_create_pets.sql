-- +goose Up
CREATE TABLE IF NOT EXISTS pets (
    id VARCHAR(36) PRIMARY KEY,
    avatar VARCHAR(10) NOT NULL DEFAULT '🐾',
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100) DEFAULT '',
    birthday VARCHAR(20) DEFAULT '',
    weight VARCHAR(20) DEFAULT '',
    notes TEXT,
    created_at BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS pets;
