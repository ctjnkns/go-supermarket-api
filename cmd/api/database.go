package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Product contains the struct details to uniquely define a product in the database
type Product struct {
	Code  string  `json:"code"`
	Name  string  `json:"name,omitempty"`
	Price float32 `json:"price,omitempty"`
}

// Products is a slice of Product structs to work with multiple items
type Products struct {
	Products []Product `json:"products"`
}

// Stringer implements for formatting printing a Product
func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

// ErrNotFound is a struct used to distinguis between errors related to badly formed json and items that were  not present in the database
type ErrNotFound struct {
	message string
}

// NewErrNotFound creates a new ErrNotFound error type
func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{
		message: message,
	}
}

// Error() returns an error of type ErrNotFound
func (e *ErrNotFound) Error() string {
	return e.message
}

// NewProduct creates a new product from individual inputs after verifying the inputs and converting to the proper format if possible
func NewProduct(code string, name string, price float32) (Product, error) {

	code, err := NewCode(code)
	if err != nil {
		return Product{}, err
	}
	name, err = NewName(name)
	if err != nil {
		return Product{}, err
	}
	err = VerifyPrice(price)
	if err != nil {
		return Product{}, err
	}

	return Product{code, name, price}, nil
}

// New Code converts the code to upper case and calls VerifyCode to check the format
func NewCode(code string) (string, error) {
	upperCode := strings.ToUpper(code)
	err := VerifyCode(upperCode)
	if err != nil {
		return "", err
	}
	return upperCode, nil
}

// NewName converts the name to lower case and calls VerifyName to check the format
func NewName(name string) (string, error) {
	lowerName := strings.ToLower(name)
	err := VerifyName(lowerName)
	if err != nil {
		return "", err
	}
	return lowerName, nil
}

// VerifyProduct checks the formatting of the product code, name and price for a specific product.
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

// VerifyCode ensures the product code is in the correct format and converts the code to upper case
func VerifyCode(code string) error {
	//upperCode := strings.ToUpper(code)
	/*
		if upperCode != code {
			return errors.New("Product Code should be uppercase")
		}
	*/

	r := regexp.MustCompile("^[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")
	if ok := r.MatchString(code); !ok {
		return fmt.Errorf("Incorrect Product Code Format: %s", code)
	}

	return nil
}

// VerifyName converts the name to lower case and makes sure that a name was provided
func VerifyName(name string) error {
	//lowerName := strings.ToLower(name)
	if name == "" {
		return errors.New("Name must be provided")
	}

	return nil
}

// VerifyPrice makes sure the price is a posive float value
func VerifyPrice(price float32) error {
	if price < 0.01 {
		return errors.New("Invalid price: Nothing is free here")
	}
	return nil
}

// InitializeDatabase adds the default products to the database
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
			return fmt.Errorf("Unable to create initial product: %s: %s\n", product, err)
		}
		_, err = app.AddDatabaseItem(p)
		if err != nil {
			return fmt.Errorf("Unable to add initial product to database: %s:  %s\n", product, err)
		}
	}
	return nil
}

// AddDatabaseItem verifies the Product fields are correct adds a single product to the database
func (app *Config) AddDatabaseItem(product Product) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	product, err := NewProduct(product.Code, product.Name, product.Price)
	if err != nil {
		return product, fmt.Errorf("Unable to create product: %w", err)
	}

	//look up the item to see if it already exists, we don't want to add duplicates.
	if _, ok := app.Database[product.Code]; ok {
		return product, fmt.Errorf("Duplicate item: %s", product.Code)
	}

	app.Database[product.Code] = product

	return product, nil
}

// AddDatabaseItems verifies takes in multiple Products, verifies each Product is valid, and adds them to product to the database
func (app *Config) AddDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products
	if len(products.Products) < 1 {
		return response, fmt.Errorf("pass in at least one product")
	}

	// look up the item to see if it already exists and verify the codes
	for _, product := range products.Products {
		product, err := NewProduct(product.Code, product.Name, product.Price)
		if err != nil {
			return response, err
		}
		if _, ok := app.Database[product.Code]; ok {
			return response, fmt.Errorf("Duplicate item: %s", product.Code)
		}
		response.Products = append(response.Products, product)
	}

	// add the items now that all have been verified
	for _, product := range response.Products {
		app.Database[product.Code] = product
	}

	return response, nil
}

// GetDatabaseItem retrieves a single product from the database, or returns an error if the product is not found or if the provided product was not valid
func (app *Config) GetDatabaseItem(code string) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	code, err := NewCode(code)
	if err != nil {
		return product, err
	}

	product, ok := app.Database[code]
	if !ok {
		return product, NewErrNotFound(fmt.Sprintf("Item not found: %s", code))
	}

	return product, nil
}

// GetDatabaseItems retrieves multiple products from the database, or returns an error if a product was invalid or not found
func (app *Config) GetDatabaseItems() (map[string]Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	products := app.Database
	if products == nil {
		return nil, errors.New("Error: empty database")
	}

	return products, nil
}

// DeleteDatabaseItem deletes a single product from the database, or returns an error if the product is not found or if the provided product was not valid
func (app *Config) DeleteDatabaseItem(code string) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	code, err := NewCode(code)
	if err != nil {
		return product, err
	}

	product, ok := app.Database[code]
	if !ok {
		return product, NewErrNotFound(fmt.Sprintf("Item not found: %s", code))
	}

	delete(app.Database, code)

	return product, nil
}

// DeleteDatabaseItems removes multiple products from the database, or returns an error if a product was invalid or not found
func (app *Config) DeleteDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products
	if len(products.Products) < 1 {
		return response, fmt.Errorf("pass in at least one product")
	}

	for _, product := range products.Products {
		code, err := NewCode(product.Code)
		if err != nil {
			return response, err
		}
		product.Code = code
		if _, ok := app.Database[product.Code]; !ok {
			return response, NewErrNotFound(fmt.Sprintf("Item not found: %s", product.Code))
		}
		response.Products = append(response.Products, product)
	}

	for _, product := range response.Products {
		delete(app.Database, product.Code)
	}

	return response, nil
}

// UpdateDatabaseItem updates a single product from the database, or returns an error if the product is not found or if the provided product was not valid
func (app *Config) UpdateDatabaseItem(product Product) (Product, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	// maybe need to update the proudct code? Assume it's an actual UUID and we don't need to
	code, err := NewCode(product.Code)
	if err != nil {
		return product, err
	}

	//look up the item to make sure it exists.
	currentProduct, ok := app.Database[code]
	if !ok {
		return currentProduct, NewErrNotFound(fmt.Sprintf("Item not found: %s", product.Code))
	}

	if product.Name != "" {
		currentProduct.Name, err = NewName(product.Name)
		if err != nil {
			return currentProduct, err
		}
	}

	if product.Price != 0 {
		err = VerifyPrice(product.Price)
		if err != nil {
			return currentProduct, err
		}
		currentProduct.Price = product.Price
	}

	app.Database[product.Code] = currentProduct

	return product, nil
}

// UpdateDatabaseItems updates multiple products in the database, or returns an error if a product was invalid or not found
func (app *Config) UpdateDatabaseItems(products Products) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var response Products

	if len(products.Products) < 1 {
		return response, fmt.Errorf("pass in at least one product")
	}

	for _, product := range products.Products {
		thisProduct, err := NewProduct(product.Code, product.Name, product.Price)
		if err != nil {
			return response, err
		}
		if _, ok := app.Database[thisProduct.Code]; !ok {
			return response, NewErrNotFound(fmt.Sprintf("Item not found: %s", product.Code))
		}
		response.Products = append(response.Products, thisProduct)
	}

	// update the items in the database now that they've all been verified
	for _, product := range response.Products {
		app.Database[product.Code] = product
	}

	return response, nil
}

// SearchDatabaseItem searches for an item by name in the database returns all the matching products or a not found error if there were no matches
func (app *Config) SearchDatabaseItem(searchString string) (Products, error) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	searchString = strings.ToLower(searchString)

	var products Products
	for _, product := range app.Database {
		if strings.Contains(product.Name, searchString) {
			products.Products = append(products.Products, product)
		}
	}
	if len(products.Products) < 1 {
		return Products{}, NewErrNotFound(fmt.Sprintf("Item not found: %s", searchString))
	}
	return products, nil
}
