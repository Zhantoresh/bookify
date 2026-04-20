INSERT INTO users (id, email, password_hash, full_name, role, phone) VALUES
('00000000-0000-0000-0000-000000000001', 'admin@booking.com', '$2a$10$oOC.td17c//5aC9FnckrVeXkmIp3d1eZbuuBGZSJciSW1BgM/I1uu', 'System Admin', 'admin', '+70000000001'),
('00000000-0000-0000-0000-000000000002', 'doctor@booking.com', '$2a$10$XrR9J7jJkCZU.tRxSR1e/.8fvI.xpCjKM8BEjuY/T5m5OfFzW53c6', 'Dr. Smith', 'provider', '+70000000002'),
('00000000-0000-0000-0000-000000000003', 'client@booking.com', '$2a$10$SZYjlq.wuqUaDxFFw.EmzeLiXx.bOF.eoMYGJHCD7QcXua.age1tu', 'John Doe', 'client', '+70000000003');

INSERT INTO services (id, provider_id, name, description, price, duration_minutes) VALUES
('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'General Consultation', 'Initial medical consultation', 100.00, 30),
('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 'Follow-up Visit', 'Follow-up medical consultation', 75.00, 20);
