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

func NewProduct(code string, name string, price float32) (Product, error) {
	code = strings.ToUpper(code)
	name = strings.ToLower(name)

	r := regexp.MustCompile("^[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")
	if ok := r.MatchString(code); !ok {
		return Product{}, errors.New("Incorrect Product Code Format")
	}

	return Product{code, name, price}, nil
}

func VerifyProduct(product Product) error {
	err := VerifyCode(product.Code)
	if err != nil {
		return err
	}
	err = VerifyName(product.Name)
	if err != nil {
		return err
	}
	err = VerifyPrice(product.Price)
	if err != nil {
		return err
	}
	return nil
}

func VerifyCode(code string) error {
	upperCode := strings.ToUpper(code)
	/*
		if upperCode != code {
			return errors.New("Product Code should be uppercase")
		}
	*/

	r := regexp.MustCompile("^[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")
	if ok := r.MatchString(upperCode); !ok {
		return errors.New("Incorrect Product Code Format")
	}

	return nil
}

func VerifyName(name string) error {
	upperName := strings.ToLower(name)
	if upperName != name {
		return errors.New("Product Name should be lowercase")
	}

	return nil
}

func VerifyPrice(price float32) error {
	if price < 0 {
		return errors.New("Nothing is free here")
	}
	return nil
}

type Config struct {
	Mutex    sync.Mutex
	Database map[string]Product
}

func (app *Config) init() {
	app.Database = make(map[string]Product)
	app.DefaultValues()
}

func (app *Config) DefaultValues() {
	initalProducts := []Product{
		{"a12T-4GH7-QPL9-3N4M", "Lettuce", 3.46}, //lower case a to test toUpper
		{"E5T6-9UI3-TH15-QR88", "Peach", 2.99},
		{"YRT6-72AS-K736-L4AR", "Green Pepper", 0.79},
		{"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", 3.59},
	}

	for _, product := range initalProducts {
		p, err := NewProduct(product.Code, product.Name, product.Price)
		if err != nil {
			log.Printf("Unable to create product: %s: %s\n", product, err)
		}
		_, err = app.AddDatabaseItem(p)
		if err != nil {
			log.Printf("Unable to add product to database: %s:  %s\n", product, err)
		}
	}
}

func (app *Config) AddDatabaseItem(product Product) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	product, err := NewProduct(product.Code, product.Name, product.Price)
	if err != nil {
		return product, errors.New(fmt.Sprintf("Unable to create product: %s", err))
	}

	//look up the item to see if it already exists, we don't want to add duplicates.
	if _, ok := app.Database[product.Code]; ok {
		return product, errors.New(fmt.Sprintf("Duplicate item: %s", product.Code))
	}

	app.Database[product.Code] = product

	return product, nil
}

func (app *Config) AddDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products

	//look up the item to see if it already exists, we don't want to add duplicates.
	for _, product := range products.Products {
		product, err := NewProduct(product.Code, product.Name, product.Price)
		if err != nil {
			return response, errors.New(fmt.Sprintf("Unable to create product: %s", err))
		}
		if _, ok := app.Database[product.Code]; ok {
			return response, errors.New(fmt.Sprintf("Duplicate item: %s", product.Code))
		}

		app.Database[product.Code] = product
		response.Products = append(response.Products, product)
	}

	return response, nil
}

func (app *Config) GetDatabaseItem(code string) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := VerifyCode(code)
	if err != nil {
		return product, errors.New(fmt.Sprintf("Unable to verify code: %s", err))
	}

	product, ok := app.Database[code]
	if !ok {
		return product, errors.New(fmt.Sprintf("Item not found: %s", code))
	}

	return product, nil
}

func (app *Config) GetDatabaseItems() (map[string]Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	products := app.Database

	return products, nil
}

func (app *Config) DeleteDatabaseItem(code string) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := VerifyCode(code)
	if err != nil {
		return product, errors.New(fmt.Sprintf("Unable to verify code: %s", err))
	}

	product, ok := app.Database[code]
	if !ok {
		return product, errors.New(fmt.Sprintf("Item not found: %s", code))
	}

	delete(app.Database, code)

	return product, nil
}

func (app *Config) DeleteDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products

	for _, product := range products.Products {
		err := VerifyCode(product.Code)
		if err != nil {
			return response, errors.New(fmt.Sprintf("Unable to verify code: %s", err))
		}
		product, ok := app.Database[product.Code]
		if !ok {
			return response, errors.New(fmt.Sprintf("Item not found: %s", product.Code))
		}

		delete(app.Database, product.Code)
		response.Products = append(response.Products, product)
	}

	return response, nil
}

func (app *Config) UpdateDatabaseItem(product Product) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	product, err := NewProduct(product.Code, product.Name, product.Price)
	if err != nil {
		return product, errors.New(fmt.Sprintf("Unable to create product: %s", err))
	}

	//look up the item to make sure it exists.
	product, ok := app.Database[product.Code]
	if !ok {
		return product, errors.New(fmt.Sprintf("Item not found: %s", product.Code))
	}

	app.Database[product.Code] = product

	return product, nil
}

func (app *Config) UpdateDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products

	for _, product := range products.Products {
		err := VerifyCode(product.Code)
		if err != nil {
			return response, errors.New(fmt.Sprintf("Unable to verify code: %s", err))
		}
		product, ok := app.Database[product.Code]
		if !ok {
			return response, errors.New(fmt.Sprintf("Item not found: %s", product.Code))
		}

		app.Database[product.Code] = product
		response.Products = append(response.Products, product)
	}

	return response, nil
}