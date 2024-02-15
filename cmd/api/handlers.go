package main

import (
	"fmt"
	"net/http"
)

func (app *Config) GetItem(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if _, ok := app.Database[product.Code]; !ok {
		msg := fmt.Sprintf("item not found %q", product.Code)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("List product: %s", product.Code),
		Data:    app.Database[product.Code],
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	payload := jsonResponse{
		Error:   false,
		Message: "Listing products",
		Data:    app.Database,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) AddItem(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//fmt.Fprintf(w, "%q\n", &product)
	//look up the item to see if it already exists, we don't want to add duplicates.
	if _, ok := app.Database[product.Code]; ok {
		msg := fmt.Sprintf("duplicate item: %q", product)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	app.Database[product.Code] = product

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Added %s", product.Code),
		Data:    product,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) AddItems(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var products Products

	err := app.readJSON(w, r, &products)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//fmt.Fprintf(w, "%q\n", &product)
	//look up the item to see if it already exists, we don't want to add duplicates.
	for _, product := range products.Products {
		if _, ok := app.Database[product.Code]; ok {
			msg := fmt.Sprintf("duplicate item: %q", product)
			http.Error(w, msg, http.StatusBadRequest) //send a 400
			return
		}

		app.Database[product.Code] = product
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Added multiple products",
		Data:    products,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) DeleteItem(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if _, ok := app.Database[product.Code]; !ok {
		msg := fmt.Sprintf("item not found %q", product.Code)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	delete(app.Database, product.Code)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Deleted product: %s", product.Code),
		Data:    nil,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) DeleteItems(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var products Products

	err := app.readJSON(w, r, &products)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	for _, product := range products.Products {
		if _, ok := app.Database[product.Code]; !ok {
			msg := fmt.Sprintf("item not found %q", product.Code)
			http.Error(w, msg, http.StatusBadRequest) //send a 400
			return
		}
		delete(app.Database, product.Code)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Deleted multiple products",
		Data:    products,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) UpdateItem(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//fmt.Fprintf(w, "%q\n", &product)
	//look up the item to see if it already exists, we don't want to add duplicates.
	if _, ok := app.Database[product.Code]; !ok {
		msg := fmt.Sprintf("item not found %q", product.Code)
		http.Error(w, msg, http.StatusBadRequest) //send a 400
		return
	}

	app.Database[product.Code] = product

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Updated %s", product.Code),
		Data:    product,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) UpdateItems(w http.ResponseWriter, r *http.Request) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	var products Products

	err := app.readJSON(w, r, &products)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//fmt.Fprintf(w, "%q\n", &product)
	//look up the item to see if it already exists, we don't want to add duplicates.
	for _, product := range products.Products {
		if _, ok := app.Database[product.Code]; !ok {
			msg := fmt.Sprintf("item not found %q", product.Code)
			http.Error(w, msg, http.StatusBadRequest) //send a 400
			return
		}

		app.Database[product.Code] = product
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Updated multiple products",
		Data:    products,
	}

	app.writeJSON(w, http.StatusOK, payload)
}
