CREATE TABLE storehouse_utility_units(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(20)
);