CREATE OR REPLACE PROCEDURE public.reserve_medicines_for_prescription(prescrip_id uuid, reason character varying)
	LANGUAGE plpgsql
AS $procedure$
	declare prescription_medicine record;
	declare stock record;
	declare pieces_needed integer;
	BEGIN
		for prescription_medicine in select id,medicine_key, m."name" as medicine_name, pieces, pieces_reserved
					from prescriptions_medicines
					left join medicines m
					on prescriptions_medicines.medicine_key = m."key" 
					where prescription_id = prescrip_id
		loop	
			
		pieces_needed := prescription_medicine.pieces - prescription_medicine.pieces_reserved;
				for stock in select * from public.get_pharmacy_stocks_sorted(prescription_medicine.medicine_key,'green') 
						union all
						select * from public.get_pharmacy_stocks_sorted(prescription_medicine.medicine_key,'ambar')
						union all
						select * from public.get_pharmacy_stocks_sorted(prescription_medicine.medicine_key,'red')
				loop
					if(pieces_needed <= 0) then
						EXIT;
					end if;
					if(stock.pieces_left = 0) then
						continue;
					end if;
					if stock.pieces_left >= pieces_needed then
						update pharmacy_stocks 
							set pieces_used = stock.pieces_used + pieces_needed
							where id = stock.id; 
						pieces_needed := 0;
					else
						update pharmacy_stocks 
							set pieces_used = stock.pieces_used + stock.pieces_left 
							where id = stock.id; 
						pieces_needed := pieces_needed - stock.pieces_left;
					end if;
				end loop;
			if(pieces_needed > 0) then
				raise exception 'No se pudo % la receta debido a la falta de piezas del medicamento: % (%)',reason,prescription_medicine.medicine_name,prescription_medicine.medicine_key;
			end if;
			update prescriptions_medicines set pieces_reserved = prescription_medicine.pieces where id = prescription_medicine.id; 		end loop;
	END;
$procedure$
;
