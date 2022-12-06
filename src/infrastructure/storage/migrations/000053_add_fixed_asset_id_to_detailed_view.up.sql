DROP VIEW fixed_assets_detailed;
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
deps.floor_number as department_floor_number,
us3.id as department_responsible_user_id,
us3.full_name as department_responsible_user_name,
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
ON fa.administrator_user_id = us2.id
INNER JOIN users us3
ON fa.administrator_user_id = us3.id;