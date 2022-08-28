CREATE TABLE prescriptions(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(9) NOT NULL,
    prescription_status_id UUID NOT NULL,
    folio SERIAL NOT NULL,
    patient_name VARCHAR(200) NOT NULL,
    observations text,
    instructions text,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    supplied_at TIMESTAMP NOT NULL,
    deletd_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_prescription_status FOREIGN KEY(prescription_status_id) REFERENCES prescriptions_statues(id)
);