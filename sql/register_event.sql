create or replace function register_event(eventName varchar(255))
returns integer as
$$
declare eventId integer;
BEGIN
insert into register_events (name) values(eventName)
returning id into eventId;
return eventId;
END;
$$ language plpgsql;