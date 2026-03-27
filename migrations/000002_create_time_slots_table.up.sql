CREATE TABLE time_slots (
    id SERIAL PRIMARY KEY,
    specialist_id INTEGER NOT NULL REFERENCES specialists(id) ON DELETE CASCADE,
    time TIMESTAMP NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_time_slots_specialist_id ON time_slots(specialist_id);

-- For John Smith (barber)
INSERT INTO time_slots (specialist_id, time) VALUES
    (1, '2026-03-27 09:00:00'),
    (1, '2026-03-27 10:00:00'),
    (1, '2026-03-27 11:00:00'),
    (1, '2026-03-28 09:00:00'),
    (1, '2026-03-28 14:00:00');

-- For Dr. Sarah Johnson (dentist)
INSERT INTO time_slots (specialist_id, time) VALUES
    (2, '2026-03-27 08:00:00'),
    (2, '2026-03-27 09:30:00'),
    (2, '2026-03-27 14:00:00'),
    (2, '2026-03-29 10:00:00');

-- For Dr. Michael Lee (psychologist)
INSERT INTO time_slots (specialist_id, time) VALUES
    (3, '2026-03-27 15:00:00'),
    (3, '2026-03-28 11:00:00'),
    (3, '2026-03-29 16:00:00'),
    (3, '2026-03-30 10:00:00');
