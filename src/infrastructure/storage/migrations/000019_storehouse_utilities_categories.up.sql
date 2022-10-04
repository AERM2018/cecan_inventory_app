CREATE TABLE storehouse_utilities_categories(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(100)
);