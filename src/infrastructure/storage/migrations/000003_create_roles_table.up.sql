CREATE TABLE roles(
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(50) NOT NULL
);