CREATE TABLE storehouse_utilities_storehouse_requests(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    storehouse_utility_key VARCHAR(15) NOT NULL,
    storehouse_request_id UUID NOT NULL,
    pieces_requested INTEGER NOT NULL,
    pieces_given INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_storehouse_utility FOREIGN KEY(storehouse_utility_key) REFERENCES storehouse_utilities(key),
    CONSTRAINT fk_storehouse_request FOREIGN KEY(storehouse_request_id) REFERENCES storehouse_requests(id)
);