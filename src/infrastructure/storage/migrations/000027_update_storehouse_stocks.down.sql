ALTER TABLE storehouse_stocks 
    DROP COLUMN quantity_presentation,
    DROP COLUMN quantity_presentation_used;

ALTER TABLE storehouse_stocks 
    RENAME COLUMN quantity_parsed TO pieces;

ALTER TABLE storehouse_stocks 
    RENAME COLUMN quantity_parsed_used TO pieces_used;

ALTER TABLE storehouse_stocks
    ALTER COLUMN pieces TYPE INT,
    ALTER COLUMN pieces_used TYPE INT;

ALTER TABLE storehouse_stocks
    ALTER COLUMN pieces_used DROP DEFAULT; 