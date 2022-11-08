CREATE OR REPLACE FUNCTION public.get_departments_ordered_by_floor(include_deleted boolean, limit_n int, offset_n int)
    RETURNS TABLE(
        id uuid,
        responsible_user_id varchar,
        name varchar,
        floor_number varchar,
        created_at timestamp,
        updated_at timestamp,
        deleted_at timestamp
    )
    LANGUAGE plpgsql
    AS $function$
        DECLARE rec record;
        BEGIN
            DROP TABLE IF EXISTS temp;
            CREATE TEMPORARY TABLE IF NOT EXISTS temp(
                id uuid,
                responsible_user_id varchar(255),
                name varchar(255),
                floor_number varchar(255),
                created_at timestamp,
                updated_at timestamp,
                deleted_at timestamp
            );
            FOR rec IN SELECT DISTINCT(d.floor_number) FROM departments d ORDER BY d.floor_number ASC
            LOOP
                IF include_deleted = FALSE THEN
                    INSERT INTO temp SELECT * FROM departments dpt WHERE dpt.floor_number = rec.floor_number AND dpt.deleted_at is null;
                ELSE
                    INSERT INTO temp SELECT * FROM departments dpt WHERE dpt.floor_number = rec.floor_number;
                END IF;
            END LOOP;
            RETURN QUERY (SELECT * FROM temp limit limit_n offset offset_n);
        END;
    $function$
;