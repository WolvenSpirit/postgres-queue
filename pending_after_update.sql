create or replace function t_notify_pending()
returns TRIGGER AS
$$
BEGIN
/*
    TODO
    Would read events for register_events table
    for each event perform notify_pending
*/
    return NEW;
END;
$$ language plpgsql;

create trigger trigger_notify_pending
after UPDATE
on notify_events
for each row
execute procedure t_notify_pending();