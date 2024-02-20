package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

const defaultPort = "8080"

// Config struct contains the database and mutex as well as other settings needed by the supermarket api
type Config struct {
	Mutex              sync.Mutex
	Database           map[string]Product
	MasJSONSize        int
	AllowUnknownFields bool
}

// init creates the database and initializes it with default data
func (app *Config) init() error {
	app.Database = make(map[string]Product)
	err := app.InitializeDatabase()
	if err != nil {
		return fmt.Errorf("Unable to initialize database with default vaules: %s", err)
	}
	return nil
}

// main initializes the database and settings and starts the web server
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("Initializing database...")
	app := Config{}
	err := app.init()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting Supermarket API on port %s\n", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
