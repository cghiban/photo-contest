
INSERT INTO auth_user (email, name, passw, created) VALUES 
	('admin@dnalc.org', 'Admin Photographer', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2021-08-13 09:04:00'),
	('user@example.com', 'User Photographer', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2021-08-13 12:20:00')
	ON CONFLICT DO NOTHING;


INSERT INTO photos (photo_id, owner_id, title, description, deleted, created_on, updated_on, updated_by) VALUES 
    ("e7b9e4e8-7b15-47ac-9d93-3d0ac42b1d46", 2, "Test photo 1", "Hopa Hopa Penelopa", false, "2021-08-18 11:59:52", "2021-08-18 11:59:52", "Admin Photographer");
