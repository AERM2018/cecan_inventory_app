CREATE TABLE medicines(
    key VARCHAR(15) PRIMARY KEY NOT NULL,
    name VARCHAR(200) NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    deleted_at DATE
);