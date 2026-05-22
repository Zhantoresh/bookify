INSERT INTO users (id, email, password_hash, full_name, role, phone) VALUES
('00000000-0000-0000-0000-000000000001', 'admin@booking.com', '$2a$10$oOC.td17c//5aC9FnckrVeXkmIp3d1eZbuuBGZSJciSW1BgM/I1uu', 'System Admin', 'admin', '+70000000001'),
ON CONFLICT (email) DO NOTHING;