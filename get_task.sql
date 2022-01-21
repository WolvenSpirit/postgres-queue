create or replace function get_task(consumerId varchar(255))
returns get_task_return as 
$func$
declare task_data get_task_return;
BEGIN
SELECT notify_events.id, notify_events.data 
INTO task_data.eventId, task_data.config 
FROM notify_events
WHERE (ack IS NULL OR ack IS FALSE)
LIMIT 1
FOR UPDATE SKIP LOCKED;
PERFORM ack(task_data.eventId, task_data.config);
return task_data;
END;
$func$ LANGUAGE plpgsql;

