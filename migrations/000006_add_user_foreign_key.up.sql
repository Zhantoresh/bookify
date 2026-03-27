ALTER TABLE bookings ADD CONSTRAINT fk_bookings_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
