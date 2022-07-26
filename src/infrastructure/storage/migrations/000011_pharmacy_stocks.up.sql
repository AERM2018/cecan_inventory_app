CREATE TABLE pharmacy_stocks(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    medicine_key VARCHAR(15) NOT NULL,
    lot_number VARCHAR(100) NOT NULL,
    pieces INT NOT NULL,
    semaforization_color SEMAFORIZATION_COLORS,
    created_at DATE NOT NULL,
    expires_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    deleted_at DATE,
    CONSTRAINT fk_medicines FOREIGN KEY(medicine_key) REFERENCES medicines(key)
);  