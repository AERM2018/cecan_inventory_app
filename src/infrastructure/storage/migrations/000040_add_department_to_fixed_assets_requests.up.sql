ALTER TABLE fixed_assets_requests ADD department_id uuid not null;
ALTER TABLE fixed_assets_requests ADD CONSTRAINT fk_department FOREIGN KEY(department_id) REFERENCES departments(id);