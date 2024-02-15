package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Code string
type dollars float32

type Product struct {
	Code  Code    `json:"code"`
	Name  string  `json:"name"`
	Price dollars `json:"price"`
}

type Products struct {
	Products []Product `json:"products"`
}

// implement Stringer for formatting prints
func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

type Config struct {
	Mutex    sync.Mutex
	Database map[Code]Product
}

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := Config{
		Database: map[Code]Product{
			"A12T-4GH7-QPL9-3N4M": {"A12T-4GH7-QPL9-3N4M", "Lettuce", 3.46},
			"E5T6-9UI3-TH15-QR88": {"E5T6-9UI3-TH15-QR88", "Peach", 2.99},
			"YRT6-72AS-K736-L4AR": {"YRT6-72AS-K736-L4AR", "Green Pepper", 0.79},
			"TQ4C-VV6T-75ZX-1RMR": {"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", 3.59},
		},
	}
	/*
		//fmt.Println(db.db["A12T-4GH7-QPL9-3N4M"])
		http.HandleFunc("/listQuery", db.listQuery)
		http.HandleFunc("/addQuery", db.addQuery)
		http.HandleFunc("/updateQuery", db.updateQuery)
		http.HandleFunc("/readQuery", db.readQuery)
		http.HandleFunc("/deleteQuery", db.deleteQuery)
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/
	log.Printf("Starting Supermarket API on port %s\n", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
