-- Drop constraints that avoid the drop of the table
ALTER TABLE fixed_assets DROP CONSTRAINT fk_description;
ALTER TABLE fixed_assets DROP CONSTRAINT fk_department;
ALTER TABLE fixed_assets DROP CONSTRAINT fk_director_user;
ALTER TABLE fixed_assets DROP CONSTRAINT fk_adminstrator_user;
ALTER TABLE fixed_assets_items_requests DROP CONSTRAINT fk_fixed_asset;
-- Drop view that is built based on this table
DROP VIEW fixed_assets_detailed;
DROP TABLE IF EXISTS fixed_assets;
-- Recreate table
CREATE TABLE fixed_assets(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    key VARCHAR(50) NOT NULL,
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
-- Recreate constraint in fixed_assets_items_requests
ALTER TABLE fixed_assets_items_requests DROP COLUMN fixed_asset_key;
ALTER TABLE fixed_assets_items_requests ADD COLUMN fixed_asset_id UUID NOT NULL;
ALTER TABLE fixed_assets_items_requests ADD CONSTRAINT fk_fixed_asset FOREIGN KEY(fixed_asset_id) REFERENCES fixed_assets(id);
-- Recreate viewt that was deleted
CREATE VIEW fixed_assets_detailed AS 
SELECT 
fa.id as id,
key,
descs.description as description,
descs.brand as brand,
descs.model as model,
series,
type,
physic_state,
deps.id as department_id,
deps.name as department_name,
observation,
us1.id as director_user_id,
us1.full_name as director_user_name,
us2.id as administrator_user_id,
us2.full_name as administrator_user_name,
fa.created_at as created_at,
fa.updated_at as updated_at
FROM fixed_assets fa
INNER JOIN fixed_asset_descriptions descs
ON fa.fixed_asset_description_id = descs.id
INNER JOIN departments deps
ON fa.department_id = deps.id
INNER JOIN users us1
ON fa.director_user_id = us1.id
INNER JOIN users us2
ON fa.administrator_user_id = us2.id;