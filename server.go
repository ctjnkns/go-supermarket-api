package main

import (
	"fmt"
	"log"
	"net/http"
)

type Code string

type Product struct {
	Name  string
	Price float32
}

// implement Stringer for formatting prints
func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

type database struct {
	products map[Code]Product
}

func (d *database) list(w http.ResponseWriter, r *http.Request) {
	for code, product := range d.products {
		fmt.Fprintf(w, "%s: %s\n", code, product)
	}
}

func main() {
	db := database{
		products: map[Code]Product{
			"A12T-4GH7-QPL9-3N4M": {"Lettuce", 3.46},
			"E5T6-9UI3-TH15-QR88": {"Peach", 2.99},
			"YRT6-72AS-K736-L4AR": {"Green Pepper", 0.79},
			"TQ4C-VV6T-75ZX-1RMR": {"Gala Apple", 3.59},
		},
	}
	//fmt.Println(db.db["A12T-4GH7-QPL9-3N4M"])
	http.HandleFunc("/list", db.list)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
