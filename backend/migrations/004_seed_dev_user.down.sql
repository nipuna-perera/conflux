-- Rollback development seed data
DELETE FROM users WHERE email = 'dev@conflux.local';
