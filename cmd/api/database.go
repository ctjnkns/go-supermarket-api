package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

	err := VerifyCode(code)
	if err != nil {
		return Product{}, fmt.Errorf("Unable to verify code: %s", err)
	}
	err = VerifyName(name)
	if err != nil {
		return Product{}, fmt.Errorf("Unable to verify name: %s", err)
	}
	err = VerifyPrice(price)
	if err != nil {
		return Product{}, fmt.Errorf("Unable to verify price: %s", err)
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
	lowerName := strings.ToLower(name)
	if lowerName == "" {
		return errors.New("Name must be provided")
	}

	return nil
}

func VerifyPrice(price float32) error {
	if price < 0 {
		return errors.New("Nothing is free here")
	}
	return nil
}

func (app *Config) InitializeDatabase() error {
	initalProducts := []Product{
		{"a12T-4GH7-QPL9-3N4M", "Lettuce", 3.46}, //lower case a to test toUpper
		{"E5T6-9UI3-TH15-QR88", "Peach", 2.99},
		{"YRT6-72AS-K736-L4AR", "Green Pepper", 0.79},
		{"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", 3.59},
	}

	for _, product := range initalProducts {
		p, err := NewProduct(product.Code, product.Name, product.Price)
		if err != nil {
			return fmt.Errorf("Unable to create product: %s: %s\n", product, err)
		}
		_, err = app.AddDatabaseItem(p)
		if err != nil {
			return fmt.Errorf("Unable to add product to database: %s:  %s\n", product, err)
		}
	}
	return nil
}

func (app *Config) AddDatabaseItem(product Product) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	product, err := NewProduct(product.Code, product.Name, product.Price)
	if err != nil {
		return product, fmt.Errorf("Unable to create product: %s", err)
	}

	//look up the item to see if it already exists, we don't want to add duplicates.
	if _, ok := app.Database[product.Code]; ok {
		return product, fmt.Errorf("Duplicate item: %s", product.Code)
	}

	app.Database[product.Code] = product

	return product, nil
}

func (app *Config) AddDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products
	if len(products.Products) < 1 {
		return response, fmt.Errorf("pass in at least one product")
	}

	//look up the item to see if it already exists, we don't want to add duplicates.
	for _, product := range products.Products {
		product, err := NewProduct(product.Code, product.Name, product.Price)
		if err != nil {
			return response, fmt.Errorf("Unable to create product: %s", err)
		}
		if _, ok := app.Database[product.Code]; ok {
			return response, fmt.Errorf("Duplicate item: %s", product.Code)
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
		return product, fmt.Errorf("Unable to verify code: %s", err)
	}

	product, ok := app.Database[code]
	if !ok {
		return product, fmt.Errorf("Item not found: %s", code)
	}

	return product, nil
}

func (app *Config) GetDatabaseItems() (map[string]Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	products := app.Database
	if products == nil {
		return nil, errors.New("Error: empty database")
	}

	return products, nil
}

func (app *Config) DeleteDatabaseItem(code string) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := VerifyCode(code)
	if err != nil {
		return product, fmt.Errorf("Unable to verify code: %s", err)
	}

	product, ok := app.Database[code]
	if !ok {
		return product, fmt.Errorf("Item not found: %s", code)
	}

	delete(app.Database, code)

	return product, nil
}

func (app *Config) DeleteDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products
	if len(products.Products) < 1 {
		return response, fmt.Errorf("pass in at least one product")
	}

	for _, product := range products.Products {
		err := VerifyCode(product.Code)
		if err != nil {
			return response, fmt.Errorf("Unable to verify code: %s", err)
		}
		product, ok := app.Database[product.Code]
		if !ok {
			return response, fmt.Errorf("Item not found: %s", product.Code)
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
		return product, fmt.Errorf("Unable to create product: %s", err)
	}

	//look up the item to make sure it exists.
	if response, ok := app.Database[product.Code]; !ok {
		return response, fmt.Errorf("Item not found: %s", product.Code)
	}

	app.Database[product.Code] = product

	return product, nil
}

func (app *Config) UpdateDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products

	if len(products.Products) < 1 {
		return response, fmt.Errorf("pass in at least one product")
	}

	for _, product := range products.Products {
		err := VerifyCode(product.Code)
		if err != nil {
			return response, fmt.Errorf("unable to verify code: %s", err)
		}
		if _, ok := app.Database[product.Code]; !ok {
			return response, fmt.Errorf("item not found: %s", product.Code)
		}

		app.Database[product.Code] = product
		response.Products = append(response.Products, product)
	}

	return response, nil
}

func (app *Config) SearchDatabaseItem(searchString string) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var products Products
	for _, product := range app.Database {
		if strings.Contains(product.Name, searchString) {
			products.Products = append(products.Products, product)
		}
	}
	if len(products.Products) < 1 {
		return Products{}, fmt.Errorf("item not found")
	}
	return products, nil
}
