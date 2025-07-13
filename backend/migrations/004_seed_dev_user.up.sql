-- Development seed data - creates default user for development environment
-- This should only run in development, not production
-- Password is 'password123' (bcrypt hashed)

INSERT INTO users (email, password_hash, first_name, last_name, created_at, updated_at)
SELECT 
    'dev@conflux.local',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt of 'password123'
    'Dev',
    'User',
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'dev@conflux.local'
);
