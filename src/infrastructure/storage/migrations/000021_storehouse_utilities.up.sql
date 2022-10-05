CREATE TABLE storehouse_utilities(
    key VARCHAR(15) PRIMARY KEY NOT NULL,
    generic_name VARCHAR(150) NOT NULL,
    presentation VARCHAR(40) NOT NULL,
    description VARCHAR(200),
    storehouse_utilities_category_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_storehouse_category FOREIGN KEY(storehouse_utilities_category_id) REFERENCES storehouse_utility_categories(id)
);