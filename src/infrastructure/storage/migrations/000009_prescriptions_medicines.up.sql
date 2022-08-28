CREATE TABLE prescriptions_medicines(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    prescription_id UUID NOT NULL,
    medicine_key VARCHAR(15) NOT NULL,
    pieces INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_prescriptions FOREIGN KEY(prescription_id) REFERENCES prescriptions(id),
    CONSTRAINT fk_medicines FOREIGN KEY(medicine_key) REFERENCES medicines(key)
);