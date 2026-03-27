CREATE TABLE specialists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO specialists (name, type) VALUES
    ('John Smith', 'barber'),
    ('Dr. Sarah Johnson', 'dentist'),
    ('Dr. Michael Lee', 'psychologist');
