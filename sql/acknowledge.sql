create or replace function ack(eventId integer, consumerId varchar(255)) RETURNS TEXT AS 
$func$
declare row_affected integer;
declare row_data text;
BEGIN
UPDATE notify_events
SET consumer = consumerId, ack = true 
WHERE id = eventId AND ack IS NOT TRUE
RETURNING ID, DATA INTO row_affected, row_data;
RETURN row_data;
END;
$func$ LANGUAGE plpgsql;