package main

import (
	"authentication/cmd/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// necessary driver for the postgres connection
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

var counts int64

func main() {
	log.Println("Starting the authentication service.")

	//Connect to the database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to postgres")
	}

	//set up config

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Check if the database is up and running
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected successfully to postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
