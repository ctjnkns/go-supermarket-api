package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) GetItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.GetDatabaseItem(product.Code)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("List product: %s", response.Code),
		Data:    response,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	response, err := app.GetDatabaseItems()
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Listing products",
		Data:    response,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) AddItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.AddDatabaseItem(product)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Added %s", product.Code),
		Data:    response,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) AddItems(w http.ResponseWriter, r *http.Request) {
	var products Products

	err := app.readJSON(w, r, &products)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.AddDatabaseItems(products)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Added multiple products",
		Data:    response,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.DeleteDatabaseItem(product.Code)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Deleted product: %s", product.Code),
		Data:    response,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) DeleteItems(w http.ResponseWriter, r *http.Request) {
	var products Products

	err := app.readJSON(w, r, &products)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.DeleteDatabaseItems(products)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Deleted multiple products",
		Data:    response,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.UpdateDatabaseItem(product)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Updated %s", product.Code),
		Data:    response,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) UpdateItems(w http.ResponseWriter, r *http.Request) {
	var products Products

	err := app.readJSON(w, r, &products)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	response, err := app.UpdateDatabaseItems(products)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("Problem with errorJSON: %s", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Updated multiple products",
		Data:    response,
	}

	app.writeJSON(w, http.StatusOK, payload)
}
