
INSERT INTO auth_user (email, name, passw, created) VALUES 
	('admin@dnalc.org', 'Admin Photographer', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2021-08-13 09:04:00'),
	('user@example.com', 'User Photographer', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2021-08-13 12:20:00')
	ON CONFLICT DO NOTHING;

