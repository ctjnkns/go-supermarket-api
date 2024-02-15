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
	Code  Code    `json:"code"`
	Name  string  `json:"name"`
	Price dollars `json:"price"`
}

// implement Stringer for formatting prints
func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

type database struct {
	mu       sync.Mutex
	products map[Code]Product
}

func (d *database) listQuery(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for code, product := range d.products {
		fmt.Fprintf(w, "%s: %s\n", code, product)
	}
}

func (d *database) addQuery(w http.ResponseWriter, req *http.Request) {
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

	d.products[code] = Product{code, name, dollars(p)} //cast the float to the dollars type

	fmt.Fprintf(w, "added %s with details %s\n", code, d.products[code]) //could return a different response code to show we created something
}

func (d *database) updateQuery(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	code := Code(req.URL.Query().Get("code"))
	name := req.URL.Query().Get("name")
	price := req.URL.Query().Get("price")

	//if the item doesn't exist, we can't update it.
	if _, ok := d.products[code]; !ok {
		msg := fmt.Sprintf("item not found: %q", code)
		http.Error(w, msg, http.StatusNotFound) //send a 404
		return
	}

	//convert the string price into a 32 bit float
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	d.products[code] = Product{code, name, dollars(p)} //cast the float to the dollars type

	fmt.Fprintf(w, "updated %s with details %s\n", code, d.products[code]) //could return a different response code to show we created something
}

func (d *database) readQuery(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	code := Code(req.URL.Query().Get("code"))

	if _, ok := d.products[code]; !ok {
		msg := fmt.Sprintf("item not found %q", code)
		http.Error(w, msg, http.StatusNotFound) //404
		return
	}

	fmt.Fprintf(w, "item %s has details %s\n", code, d.products[code])
}

func (d *database) deleteQuery(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	code := Code(req.URL.Query().Get("code"))

	if _, ok := d.products[code]; !ok {
		msg := fmt.Sprintf("item not found %q", code)
		http.Error(w, msg, http.StatusNotFound) //404
		return
	}

	delete(d.products, code)
	fmt.Fprintf(w, "deleted %s\n", code)
}

func main() {
	db := database{
		products: map[Code]Product{
			"A12T-4GH7-QPL9-3N4M": {"A12T-4GH7-QPL9-3N4M", "Lettuce", 3.46},
			"E5T6-9UI3-TH15-QR88": {"E5T6-9UI3-TH15-QR88", "Peach", 2.99},
			"YRT6-72AS-K736-L4AR": {"YRT6-72AS-K736-L4AR", "Green Pepper", 0.79},
			"TQ4C-VV6T-75ZX-1RMR": {"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", 3.59},
		},
	}
	//fmt.Println(db.db["A12T-4GH7-QPL9-3N4M"])
	http.HandleFunc("/listQuery", db.listQuery)
	http.HandleFunc("/addQuery", db.addQuery)
	http.HandleFunc("/updateQuery", db.updateQuery)
	http.HandleFunc("/readQuery", db.readQuery)
	http.HandleFunc("/deleteQuery", db.deleteQuery)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
