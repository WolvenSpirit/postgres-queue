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

	_ "github.com/lib/pq"
)

var (
	addr      string
	dbUser    string
	dbPass    string
	dbHost    string
	dbName    string
	dbsslMode string
	DB        *sql.DB
)

func connectDB(user, pass, host, database, sslMode string) {
	var err error
	driverName := "postgres"
	url := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", driverName, user, pass, host, database, sslMode)
	if DB, err = sql.Open(driverName, url); err != nil {
		log.Println("sql.Open: ", err.Error())
	}
	if err = DB.Ping(); err != nil {
		log.Println("Failed to establish database connection")
	}
	fmt.Println("open postgres connections:", DB.Stats().OpenConnections)
}

func defineFlags() {
	flag.StringVar(&addr, "address", "0.0.0.0:9001", "Server address")
	flag.StringVar(&dbUser, "dbuser", "postgres", "database user")
	flag.StringVar(&dbPass, "dbpass", "", "database password")
	flag.StringVar(&dbHost, "dbhost", "localhost:5432", "database host")
	flag.StringVar(&dbName, "dbname", "shared_db01", "database name")
	flag.StringVar(&dbsslMode, "dbsslmode", "require", "database sslMode")
}

func listen(addr string, mux *http.ServeMux) *http.Server {
	s := http.Server{Addr: addr, Handler: mux}

	go func() {
		s.ListenAndServe()
	}()
	fmt.Println("Started listening on ", addr)
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
	defineFlags()
	flag.Parse()
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
}
