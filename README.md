# postgres-queue

## Minimal general purpose queue implementation using PostgreSQL.

---

How is it supposed to work?

Dispatcher calls `publish(dispatcherId, eventChannelName, configData)`.

Any listener on that namespace will be notified, the first one that replies will claim the task.

A cron job or routine will check the table for rows that are not acknowledged.

If there are such rows it will broadcast them again.

Any new worker/consumer that just spawned and thus did not listen or receive any task will call `get_task(consumerId)` to check for pending tasks.

Consumer should listen after spawning on the designated namespace (`eventChannelName`) for pending tasks.
