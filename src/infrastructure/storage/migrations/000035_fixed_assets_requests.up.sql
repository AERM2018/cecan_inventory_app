CREATE TABLE fixed_assets_requests(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(9) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_creator_user FOREIGN KEY(user_id) REFERENCES users(id)
);