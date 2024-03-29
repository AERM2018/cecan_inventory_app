CREATE TABLE users(
  id VARCHAR(9) PRIMARY KEY NOT NULL,
  role_id UUID NOT NULL,
  name VARCHAR(100) NOT NULL,
  surname VARCHAR(100) NOT NULL,
  email VARCHAR(150) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_roles FOREIGN KEY(role_id) REFERENCES roles(id)
);