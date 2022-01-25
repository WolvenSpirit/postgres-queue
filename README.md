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

## Make it your own

By editing `tasks.go` and adding a function of type task `Task`
```go 
type Task func(payload string, eventId string, status *chan int)
```
and mapping this `Task` to a channel (this represents the pg_Notify channel name)
```go
sched.DefinedTasks["Channel01"] = SomeTask // type Task
```

## Build it
`go build -o task_runner .`
## Run it
`./task_runner`

## Send events for processing tasks defined by you

By using the publish function in postgres
```sql
select publish('dispatcher01','Channel01','{"config_prop":"some-value"}');
```
The task you defined will receive the config object stringified.

Write basic Go, define a struct and `json.Unmarshal` the payload parameter to gain access to the task parameter data sent via postgres.

A `Task` should flag when it is starting and ending.

```go
func task(payload string, eventId string, status *chan int) {
    (*status) <-StatusOngoing // This ensures that the task is properly accounted for by the task scheduler
    // Do the task logic
    (*status) <-StatusEnded // Will signal that the goroutine will exit, possibly making room for a standby task to start
}
```
