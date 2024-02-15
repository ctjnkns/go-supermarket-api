package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

type Product struct {
	Code  string  `json:"code"`
	Name  string  `json:"name,omitempty"`
	Price float32 `json:"price,omitempty"`
}

type Products struct {
	Products []Product `json:"products"`
}

// implement Stringer for formatting prints
func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

func newProduct(code string, name string, price float32) (Product, error) {
	code = strings.ToUpper(code)
	name = strings.ToLower(name)

	r := regexp.MustCompile("^[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")
	if ok := r.MatchString(code); !ok {
		return Product{}, errors.New("Incorrect Product Code Format")
	}

	return Product{code, name, price}, nil
}

type Config struct {
	Mutex    sync.Mutex
	Database map[string]Product
}

func (app *Config) init() {
	app.Database = make(map[string]Product)
	app.DefaultValues()
}

func (app *Config) Add(products ...Product) error {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	if len(products) < 1 {
		return errors.New("Must specify at least one product to add")
	}
	for _, product := range products {
		app.Database[product.Code] = product
	}
	return nil
}

func (app *Config) DefaultValues() {
	initalProducts := []Product{
		Product{"a12T-4GH7-QPL9-3N4M", "Lettuce", 3.46}, //lower case a to test toUpper
		Product{"E5T6-9UI3-TH15-QR88", "Peach", 2.99},
		Product{"YRT6-72AS-K736-L4AR", "Green Pepper", 0.79},
		Product{"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", 3.59},
	}

	for _, product := range initalProducts {
		p, err := newProduct(product.Code, product.Name, product.Price)
		if err != nil {
			log.Printf("Unable to create product: %s: %s\n", product, err)
		}
		err = app.Add(p)
		if err != nil {
			log.Printf("Unable to add product to database: %s:  %s\n", product, err)
		}
	}
}
