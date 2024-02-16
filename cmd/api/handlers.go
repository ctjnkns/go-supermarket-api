package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) GetItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.ReadJSON(w, r, &product)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	response, err := app.GetDatabaseItem(product.Code)
	if err != nil {
		err := app.ErrorJSON(w, err, http.StatusNotFound)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("List product: %s", response.Code),
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	response, err := app.GetDatabaseItems()
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Listing products",
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) AddItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.ReadJSON(w, r, &product)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	response, err := app.AddDatabaseItem(product)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Added %s", product.Code),
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) AddItems(w http.ResponseWriter, r *http.Request) {
	var products Products

	err := app.ReadJSON(w, r, &products)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	response, err := app.AddDatabaseItems(products)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Added multiple products",
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.ReadJSON(w, r, &product)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	response, err := app.DeleteDatabaseItem(product.Code)
	if err != nil {
		err := app.ErrorJSON(w, err, http.StatusNotFound)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Deleted product: %s", product.Code),
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) DeleteItems(w http.ResponseWriter, r *http.Request) {
	var products Products

	err := app.ReadJSON(w, r, &products)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	response, err := app.DeleteDatabaseItems(products)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Deleted multiple products",
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := app.ReadJSON(w, r, &product)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	response, err := app.UpdateDatabaseItem(product)
	if err != nil {
		err := app.ErrorJSON(w, err, http.StatusNotFound)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Updated %s", product.Code),
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}

func (app *Config) UpdateItems(w http.ResponseWriter, r *http.Request) {
	var products Products

	err := app.ReadJSON(w, r, &products)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}
	response, err := app.UpdateDatabaseItems(products)
	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			log.Printf("Problem with ErrorJSON: %s", err)
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Updated multiple products",
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}
