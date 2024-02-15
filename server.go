package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Code string
type dollars float32

type Product struct {
	Name  string
	Price dollars
}

// implement Stringer for formatting prints
func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

type database struct {
	mu       sync.Mutex
	products map[Code]Product
}

func (d *database) list(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for code, product := range d.products {
		fmt.Fprintf(w, "%s: %s\n", code, product)
	}
}

func (d *database) add(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	code := Code(req.URL.Query().Get("code"))
	name := req.URL.Query().Get("name")
	price := req.URL.Query().Get("price")

	//look up the item to see if it already exists, we don't want to add duplicates.
	if _, ok := d.products[code]; ok {
		msg := fmt.Sprintf("duplicate item: %q", code)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	//convert the string price into a 32 bit float
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	d.products[code] = Product{name, dollars(p)} //cast the float to the dollars type

	fmt.Fprintf(w, "added %s with details %s\n", code, d.products[code]) //could return a different response code to show we created something
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
	http.HandleFunc("/add", db.add)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
