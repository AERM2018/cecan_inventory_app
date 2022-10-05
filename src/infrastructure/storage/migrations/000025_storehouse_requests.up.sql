CREATE TABLE storehouse_requests(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(9) NOT NULL,
    folio SERIAL NOT NULL,
    storehouse_requests_status_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);