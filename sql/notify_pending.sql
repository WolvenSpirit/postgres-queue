create or replace function notify_pending(ch varchar(255)) 
returns SETOF integer as
$$
declare eventId integer;
BEGIN
    FOR eventId IN
    SELECT id FROM notify_events WHERE ack IS NULL OR ack IS FALSE
    LOOP
        RETURN NEXT eventId;
        PERFORM pg_notify(ch, eventId::text);
    END LOOP;
    RETURN;
END;
$$ LANGUAGE plpgsql;