CREATE TABLE storehouse_requests_statues(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(20) NOT NULL
);