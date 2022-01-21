create or replace function notify_pending(ch varchar(255)) 
returns void as
$$
declare eventId integer;
BEGIN
FOR COUNTER IN 1..100 LOOP
SELECT id INTO eventId FROM notify_events WHERE ack IS NULL OR ack IS FALSE LIMIT 1;
EXIT WHEN eventId IS NULL;
EXIT WHEN eventId = 0;
PERFORM pg_notify(ch, eventId::text);
END LOOP;
END;
$$ LANGUAGE plpgsql;