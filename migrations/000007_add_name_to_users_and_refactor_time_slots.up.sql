-- Add name column to users table
ALTER TABLE users ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT 'Unknown';

-- Drop specialists table since specialists are now users with role 'specialist'
DROP TABLE IF EXISTS specialists CASCADE;

-- Recreate time_slots table with user_id instead of specialist_id
-- First, backup the old data if needed, then drop and recreate
DROP TABLE IF EXISTS time_slots CASCADE;

CREATE TABLE time_slots (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    time TIMESTAMP NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_time_slots_user_id ON time_slots(user_id);
CREATE INDEX idx_time_slots_time ON time_slots(time);
