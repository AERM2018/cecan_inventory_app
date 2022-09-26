CREATE OR REPLACE FUNCTION public.get_pharmacy_stocks_sorted(med_key character varying, color character varying)
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
		return query (select
			phs.id,
			phs.pieces,
			phs.pieces_used,
			-phs.pieces_used + phs.pieces as pieces_left,
			phs.expires_at,
			phs.semaforization_color
		from pharmacy_stocks phs
		where phs.medicine_key = med_key and phs.semaforization_color::text = color) order by expires_at asc;
	END;
$function$
;
