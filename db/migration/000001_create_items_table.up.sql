/*
 I dont know what happen
 */
CREATE TABLE IF NOT EXISTS data
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    status TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);