package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	addr         string
	dbUser       string
	dbPass       string
	dbHost       string
	dbName       string
	dbsslMode    string
	DB           *sql.DB
	listener     *pq.Listener
	eventsFlag   string
	events       []string
	LStdout      *log.Logger
	LStderr      *log.Logger
	consumerName string
	sched        *Scheduler
)

const (
	Channel01 = "Channel01"
	Channel02 = "Channel02"
	Channel03 = "Channel03"
	Channel04 = "Channel04"
	Channel05 = "Channel05"
	Channel06 = "Channel06"
	Channel07 = "Channel07"
	Channel08 = "Channel08"
	Channel09 = "Channel09"
)

func loggerInit() {
	LStdout = log.New(os.Stdout, "\033[32m ", log.Ldate|log.Ltime)
	LStderr = log.New(os.Stderr, "\033[31m ", log.Ldate|log.Ltime)
}

func listenNotifyEvents(dsn string, consumerName string) {
	listener = pq.NewListener(dsn, time.Second*1, time.Second*120, func(event pq.ListenerEventType, err error) {
		if event == pq.ListenerEventConnected {
			LStdout.Println("pq listener connected")
		}
		if event == pq.ListenerEventConnectionAttemptFailed {
			LStderr.Println("pq listener connection attempt failed")
		}
		if event == pq.ListenerEventDisconnected {
			LStdout.Println("pq listener disconnected")
		}
		if event == pq.ListenerEventReconnected {
			LStdout.Println("pq listener reconnected")
		}
		if err != nil {
			LStderr.Println(err.Error())
		}
	})
	for k := range events {
		if err := listener.Listen(events[k]); err != nil {
			LStderr.Println(err.Error())
		}
	}
	go func() {
		for {
			select {
			case b := <-listener.Notify:
				LStdout.Printf("channel [%s]: %s", b.Channel, b.Extra)
				row := DB.QueryRow(fmt.Sprintf("select ack(%s,'%s');", b.Extra, consumerName))
				data := new(interface{})
				if err := row.Scan(&data); err != nil {
					LStderr.Println("scan: ", err.Error())
				}
				// TODO spawn task with deadline set depending on channel
				LStdout.Println(((*data).(string)))
				go sched.NewTask((*data).(string), b.Extra, b.Channel)
			default:
				time.Sleep(time.Second * 1)
			}
		}
	}()
}

func connectDB(user, pass, host, database, sslMode string) {
	var err error
	driverName := "postgres"
	url := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", driverName, user, pass, host, database, sslMode)
	if DB, err = sql.Open(driverName, url); err != nil {
		LStderr.Println("sql.Open: ", err.Error())
	}
	if err = DB.Ping(); err != nil {
		LStderr.Println("Failed to establish database connection")
	}
	fmt.Println("open postgres connections:", DB.Stats().OpenConnections)
	listenNotifyEvents(url, consumerName)
}

func defineFlags() {
	flag.StringVar(&addr, "address", "0.0.0.0:9001", "Server address")
	flag.StringVar(&dbUser, "dbuser", "postgres", "database user")
	flag.StringVar(&dbPass, "dbpass", "", "database password")
	flag.StringVar(&dbHost, "dbhost", "localhost:5432", "database host")
	flag.StringVar(&dbName, "dbname", "shared_db01", "database name")
	flag.StringVar(&dbsslMode, "dbsslmode", "require", "database sslMode")
	flag.StringVar(&eventsFlag, "notify", "basic,fast", "pg_notify event namespace")
	flag.StringVar(&consumerName, "consumer-name", "task-spawner", "consumer name that acks and spawns tasks")
}

func listen(addr string, mux *http.ServeMux) *http.Server {
	s := http.Server{Addr: addr, Handler: mux}

	go func() {
		s.ListenAndServe()
	}()
	LStdout.Println("Started listening on ", addr)
	return &s
}

// defineHandlers contains the mapping of dispatcher endpoints that would trigger a new event
func defineHandlers(handlers *map[string]http.HandlerFunc) {

	(*handlers)["/"] = func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		rw.Write([]byte("..."))
	}
}

func registerHandlers(mux *http.ServeMux, handlers map[string]http.HandlerFunc) {
	for k := range handlers {
		mux.HandleFunc(k, handlers[k])
	}
}

func main() {
	sched = &Scheduler{MaxConcurrentTasks: 9, TaskDeadline: 120, RetryAfter: time.Second * 9}
	sched.DefinedTasks = make(map[string]Task)
	MapTask()
	loggerInit()
	defineFlags()
	flag.Parse()
	events = strings.Split(eventsFlag, ",")
	events = append(events, []string{Channel01, Channel02, Channel03, Channel04, Channel05, Channel06, Channel07, Channel08, Channel09}...)
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	mux := new(http.ServeMux)
	handlers := make(map[string]http.HandlerFunc)
	defineHandlers(&handlers)
	registerHandlers(mux, handlers)
	connectDB(dbUser, dbPass, dbHost, dbName, dbsslMode)
	s := listen(addr, mux)
	<-sigInt
	s.Shutdown(context.Background())
	listener.Close()
	DB.Close()
}
