ALTER TABLE storehouse_stocks
    ADD COLUMN lot_number VARCHAR(15) NOT NULL,
    ADD COLUMN catalog_number VARCHAR(20) NOT NULL,
    ADD COLUMN expires_at TIMESTAMP NOT NULL;
