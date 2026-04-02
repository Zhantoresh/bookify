BEGIN;


CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 2. Добавляем админа.

INSERT INTO users (email, password_hash, role, name, created_at)
VALUES (
           'admin@bookify.kz',
           '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- это хеш пароля 'admin123'
           'admin',
           'Super Admin',
           NOW()
       )
    ON CONFLICT (email) DO UPDATE SET role = 'admin';

COMMIT;