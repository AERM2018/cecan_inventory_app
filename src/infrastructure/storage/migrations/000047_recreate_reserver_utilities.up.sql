DROP PROCEDURE IF EXISTS public.reserve_utility_to_request(IN req_id uuid, IN uti_key character varying, IN pieces_to_reserve integer);
CREATE OR REPLACE PROCEDURE public.reserve_utility_to_request(IN req_id uuid)
 LANGUAGE plpgsql
AS $procedure$
	declare storehouse_utility_storehouse_request record;
	declare total_quantity_parsed_left float;
	declare stock record;
    declare pieces_needed integer;
	BEGIN
			-- Find details about storehouse utility
			for storehouse_utility_storehouse_request in select 
				id,storehouse_utility_key,stouti."generic_name" as utility_name, stouti."quantity_per_unit" as quantity_per_unit, pieces, last_pieces_supplied
				from storehouse_utilities_storehouse_requests utireq
				left join storehouse_utilities stouti
				on utireq.storehouse_utility_key = stouti.key
				where storehouse_request_id = req_id
			loop
				pieces_needed = storehouse_utility_storehouse_request.last_pieces_supplied;
				-- Searching for stocks to complete the request utilities' quantity
				for stock in select id, quantity_parsed_used,quantity_parsed, quantity_parsed_left, generic_name 
                from get_storehouse_stocks_sorted(storehouse_utility_storehouse_request.storehouse_utility_key)
				loop
					if(pieces_needed <= 0) then
						EXIT;
					end if;
					if stock.quantity_parsed_left >= pieces_needed then
						update storehouse_stocks 
							set quantity_presentation_used = (stock.quantity_parsed_used + stock.quantity_parsed_left) / storehouse_utility_storehouse_request.quantity_per_unit,
								quantity_parsed_used = stock.quantity_parsed_used + pieces_needed
							where id = stock.id; 
						pieces_needed := 0;
					else
						update storehouse_stocks 
							set quantity_presentation_used = (stock.quantity_parsed_used + stock.quantity_parsed_left) / storehouse_utility_storehouse_request.quantity_per_unit,
								quantity_parsed_used = stock.quantity_parsed_used + stock.quantity_parsed_left 
							where id = stock.id; 
						pieces_needed := pieces_needed - stock.quantity_parsed_left;
					end if;
				end loop;
			end loop;
	END;
$procedure$
;
