DROP FUNCTION IF EXISTS public.get_pharmacy_stocks_sorted_no_color(med_key character varying);
CREATE OR REPLACE FUNCTION public.get_pharmacy_stocks_sorted_no_color(med_key character varying)
	RETURNS table (
		id uuid,
		pieces integer,
		pieces_used integer,
		pieces_left integer, 
		expires_at timestamp, 
		semaforization_color public.semaforization_colors)
	LANGUAGE plpgsql
AS $function$
	BEGIN
		return query (select * from ( select
				phs.id,
				phs.pieces,
				phs.pieces_used,
				-phs.pieces_used + phs.pieces as pieces_left,
				phs.expires_at,
				phs.semaforization_color
			from pharmacy_stocks phs
			where (phs.medicine_key = med_key)
			order by expires_at asc, created_at asc) result
		where result.pieces_left > 0)
		;
	END;
$function$
;
