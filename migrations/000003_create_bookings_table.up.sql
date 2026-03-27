CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    time_slot_id INTEGER NOT NULL REFERENCES time_slots(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'BOOKED'
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_time_slot_id ON bookings(time_slot_id);
CREATE UNIQUE INDEX idx_bookings_time_slot_booked ON bookings(time_slot_id) WHERE status = 'BOOKED';
