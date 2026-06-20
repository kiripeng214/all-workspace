CREATE TABLE IF NOT EXISTS feeding_records (
    id VARCHAR(36) PRIMARY KEY,
    pet_id VARCHAR(36) NOT NULL,
    schedule_id VARCHAR(36) DEFAULT NULL,
    time VARCHAR(10) DEFAULT '',
    food_type VARCHAR(50) DEFAULT '粮食',
    amount VARCHAR(50) DEFAULT '一份',
    notes TEXT,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
