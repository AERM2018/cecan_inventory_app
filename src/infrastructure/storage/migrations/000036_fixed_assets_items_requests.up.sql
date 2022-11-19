CREATE TABLE fixed_assets_items_requets(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fixed_asset_key VARCHAR(50) NOT NULL,
    fixed_assets_request_id UUID NOT NULL,
    CONSTRAINT fk_fixed_asset FOREIGN KEY(fixed_asset_key) REFERENCES fixed_assets(key),
    CONSTRAINT fk_fixed_assets_request FOREIGN KEY(fixed_assets_request_id) REFERENCES fixed_assets_requests(id)
);