CREATE TABLE departments(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    responsible_user_id VARCHAR(9),
    name VARCHAR(50) NOT NULL,
    floor_number VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_responsible_user FOREIGN KEY(responsible_user_id) REFERENCES users(id)
);