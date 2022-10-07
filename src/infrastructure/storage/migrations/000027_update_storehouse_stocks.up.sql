ALTER TABLE storehouse_stocks 
    ADD COLUMN quantity_presentation FLOAT NOT NULL,
    ADD COLUMN quantity_presentation_used FLOAT;

ALTER TABLE storehouse_stocks 
    RENAME COLUMN pieces TO quantity_parsed;

ALTER TABLE storehouse_stocks 
    RENAME COLUMN pieces_used TO quantity_parsed_used;

ALTER TABLE storehouse_stocks
    ALTER COLUMN quantity_parsed TYPE FLOAT,
    ALTER COLUMN quantity_parsed_used TYPE FLOAT; 

ALTER TABLE storehouse_stocks
    ALTER COLUMN quantity_presentation_used SET DEFAULT 0,
    ALTER COLUMN quantity_parsed_used SET DEFAULT 0; 