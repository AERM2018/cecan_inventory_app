CREATE TABLE storehouse_stocks(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    storehouse_utility_key VARCHAR(15) NOT NULL,
    pieces INTEGER NOT NULL,
    pieces_used INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_storehouse_utility FOREIGN KEY(storehouse_utility_key) REFERENCES storehouse_utilities(key)
);