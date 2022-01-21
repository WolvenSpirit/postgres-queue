create or replace function ack(eventId integer, consumerId varchar(255)) RETURNS INTEGER AS 
$func$
declare row_affected integer;
BEGIN
UPDATE notify_events
SET consumer = consumerId, ack = true 
WHERE id = eventId
RETURNING ID INTO row_affected;
RETURN row_affected;
END;
$func$ LANGUAGE plpgsql;