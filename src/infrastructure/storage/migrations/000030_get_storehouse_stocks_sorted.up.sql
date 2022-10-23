CREATE OR REPLACE FUNCTION public.get_storehouse_stocks_sorted(uti_key character varying)
	RETURNS table (
		id uuid,
		generic_name varchar,
		storehouse_utility_key varchar,
		quantity_parsed float,
		quantity_parsed_used float,
		quantity_parsed_left float, 
		quantity_presentation float,
		quantity_presentation_used float,
		quantity_presentation_left float, 
		expires_at timestamp)
	LANGUAGE plpgsql
AS $function$
	BEGIN
		return query select * from (select
				shs.id,
				su.generic_name,
				shs.storehouse_utility_key,
				shs.quantity_parsed,
				shs.quantity_parsed_used,
				-shs.quantity_parsed_used + shs.quantity_parsed as quantity_parsed_left,
                shs.quantity_presentation,
				shs.quantity_presentation_used,
				-shs.quantity_presentation_used + shs.quantity_presentation as quantity_presentation_left,
				shs.expires_at
			from storehouse_stocks shs
			left join storehouse_utilities su
			on su.key = shs.storehouse_utility_key
			order by shs.expires_at asc, shs.created_at asc) result
			where result.quantity_parsed_left > 0;
	END;
$function$
;
