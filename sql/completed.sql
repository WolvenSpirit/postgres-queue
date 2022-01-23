create or replace function set_completed(eventId integer, task_duration numeric(9,2))
returns void as
$$
BEGIN
update notify_events set duration = task_duration, completed = NOW()
where id = eventId;
END;
$$ language plpgsql;
