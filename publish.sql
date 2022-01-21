create or replace function publish(pub varchar(255), ch varchar(255), config text)
returns integer as $func$
declare notify_event_id integer;
BEGIN
INSERT INTO notify_events
(publisher, channel, data) values (pub, ch, config)
RETURNING ID INTO notify_event_id;
PERFORM pg_notify(ch, notify_event_id::text);
RETURN notify_event_id;
END;
$func$ LANGUAGE plpgsql;
