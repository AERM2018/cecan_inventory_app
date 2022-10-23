CREATE TABLE storehouse_utilities(
    key VARCHAR(15) PRIMARY KEY NOT NULL,
    generic_name VARCHAR(150) NOT NULL,
    storehouse_utility_presentation_id UUID NOT NULL,
    storehouse_utility_unit_id UUID NOT NULL,
    quantity_per_unit FLOAT(4) NOT NULL,
    description VARCHAR(200),
    storehouse_utility_category_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_storehouse_category FOREIGN KEY(storehouse_utility_category_id) REFERENCES storehouse_utility_categories(id),
    CONSTRAINT fk_storehouse_presentation FOREIGN KEY(storehouse_utility_presentation_id) REFERENCES storehouse_utility_presentations(id),
    CONSTRAINT fk_storehouse_unit FOREIGN KEY(storehouse_utility_unit_id) REFERENCES storehouse_utility_units(id)
);