-- Revert: Drop time_slots and recreate specialists table
DROP TABLE IF EXISTS time_slots;

-- Recreate specialists table
CREATE TABLE specialists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL
);

-- Recreate time_slots with specialist_id
CREATE TABLE time_slots (
    id SERIAL PRIMARY KEY,
    specialist_id INTEGER NOT NULL REFERENCES specialists(id) ON DELETE CASCADE,
    time TIMESTAMP NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_time_slots_specialist_id ON time_slots(specialist_id);

-- Remove name column from users
ALTER TABLE users DROP COLUMN name;
