DROP PROCEDURE IF EXISTS public.reserve_utility_to_request(IN req_id uuid);
CREATE OR REPLACE PROCEDURE public.reserve_utility_to_request(IN req_id uuid, IN uti_key character varying, IN pieces_to_reserve integer)
 LANGUAGE plpgsql
AS $procedure$
	declare request_utility record;
	declare total_quantity_parsed_left float;
	declare stock record;
    declare pieces_needed integer;
	BEGIN
	 			pieces_needed := pieces_to_reserve;
				-- Find details about storehouse utility
				select generic_name, shu.name as unit into request_utility 
				from storehouse_utilities 
				left join storehouse_utility_units shu
				on shu.id = storehouse_utilities.storehouse_utility_unit_id
				where key = uti_key;
				-- Not enough pieces in stocks
				select sum(quantity_parsed_left) into total_quantity_parsed_left from get_storehouse_stocks_sorted(uti_key) stocks group by stocks.storehouse_utility_key;
				if pieces_needed > total_quantity_parsed_left then
					RAISE EXCEPTION 'No hay suficiente stock para suministrar el elemento de almacen %, % % disponibles.', request_utility.generic_name,total_quantity_parsed_left,lower(request_utility.unit);
				end if;
				-- Searching for stocks to complete the request utilities' quantity
				for stock in select id, quantity_parsed_used,quantity_parsed, quantity_parsed_left, generic_name 
                from get_storehouse_stocks_sorted(uti_key)
				loop
					if(pieces_needed <= 0) then
						EXIT;
					end if;
					if stock.quantity_parsed_left >= pieces_needed then
						update storehouse_stocks 
							set quantity_parsed_used = stock.quantity_parsed_used + pieces_needed
							where id = stock.id; 
						pieces_needed := 0;
					else
						update storehouse_stocks 
							set quantity_parsed_used = stock.quantity_parsed_used + stock.quantity_parsed_left 
							where id = stock.id; 
						pieces_needed := pieces_needed - stock.quantity_parsed_left;
					end if;
				end loop;
				
				if pieces_needed != pieces_to_reserve then
                	-- Update the storehouse utility reference to the request
                	update storehouse_utilities_storehouse_requests set pieces_supplied = pieces_supplied + pieces_to_reserve
                	where(storehouse_request_id = req_id AND storehouse_utility_key = uti_key);
				else
					-- When there are 0 pieces in stock
					RAISE EXCEPTION 'No hay stock para suministrar el elemento de almacen %', request_utility.generic_name;
				end if;

	END;
$procedure$
;
