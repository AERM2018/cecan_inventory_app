ALTER TABLE fixed_assets
ADD department_responsible_user_id VARCHAR(9);
ALTER TABLE fixed_assets
ADD CONSTRAINT fk_department_responsible FOREIGN KEY(department_responsible_user_id) REFERENCES users(id);
UPDATE fixed_assets SET department_responsible_user_id = (SELECT id FROM users LIMIT 1);
ALTER TABLE fixed_assets ALTER COLUMN department_responsible_user_id SET NOT NULL;