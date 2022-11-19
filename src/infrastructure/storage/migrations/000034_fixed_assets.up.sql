CREATE TABLE fixed_assets(
    key VARCHAR(50) PRIMARY KEY,
    fixed_asset_description_id UUID NOT NULL,
    series VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    physic_state VARCHAR(50) NOT NULL,
    department_id UUID NOT NULL,
    observation VARCHAR(255),
    director_user_id VARCHAR(9) NOT NULL,
    administrator_user_id VARCHAR(9) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_description FOREIGN KEY(fixed_asset_description_id) REFERENCES fixed_asset_descriptions(id),
    CONSTRAINT fk_department FOREIGN KEY(department_id) REFERENCES departments(id),
    CONSTRAINT fk_director_user FOREIGN KEY(director_user_id) REFERENCES users(id),
    CONSTRAINT fk_adminstrator_user FOREIGN KEY(administrator_user_id) REFERENCES users(id)
);