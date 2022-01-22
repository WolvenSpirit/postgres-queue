create or replace function notify_pending_wrapper()
returns SETOF text as
$$
declare event_name text;
BEGIN
    FOR event_name IN
    SELECT name from register_events re WHERE re.skip IS FALSE
    LOOP
        RETURN NEXT event_name;
        SELECT notify_pending(event_name);
    END LOOP;   
    RETURN;
END;
$$ language plpgsql;

create or replace function t_notify_pending()
returns TRIGGER AS
$$
BEGIN
    PERFORM notify_pending_wrapper();
    RETURN NEW;
END;
$$ language plpgsql;

create trigger trigger_notify_pending
after UPDATE
on notify_events
for each row
EXECUTE PROCEDURE t_notify_pending();