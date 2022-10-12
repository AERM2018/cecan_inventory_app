CREATE TABLE storehouse_requests(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(9) NOT NULL,
    folio SERIAL NOT NULL,
    storehouse_request_status_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_request_status FOREIGN KEY(storehouse_request_status_id) REFERENCES storehouse_request_statuses(id),
    CONSTRAINT fk_request_user FOREIGN KEY(user_id) REFERENCES users(id)
);