create or replace function notify_pending(ch varchar(255)) 
returns varchar(255) as
$$
declare eventId integer;
declare i integer = 0;
declare return_message varchar(255);
BEGIN
FOR COUNTER IN 1..100 LOOP
SELECT id INTO eventId FROM notify_events WHERE ack IS NULL OR ack IS FALSE LIMIT 1;
EXIT WHEN eventId IS NULL;
EXIT WHEN eventId = 0;
PERFORM pg_notify(ch, eventId::text);
i = i + 1;
END LOOP;
SELECT 'pg_notify performed ' || i::varchar(255) || ' times.' as msg into return_message;
RETURN return_message;
END;
$$ LANGUAGE plpgsql;