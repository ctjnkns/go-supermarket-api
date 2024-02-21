package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// GetItem retrieves a sinlge item from the database and writes it to the http response writer, or returns an error if the item was not found
func (app *Config) GetItem(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	response, err := app.GetDatabaseItem(code)
	if err != nil {
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// GetItems retrieves all items from the database and writes them to the http response writer
func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	response, err := app.GetDatabaseItems()
	if err != nil {
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// AddItem adds a sinlge item to the database and writes a response to the http response writer
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
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// AddItems adds multiple items to the database and writes writes a response to the http response writer
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
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// DeleteItem deletes a sinlge item from the database and writes the details to the http response writer, or returns an error if the item was not found
func (app *Config) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var product Product

	product.Code = chi.URLParam(r, "code")

	response, err := app.DeleteDatabaseItem(product.Code)
	if err != nil {
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// DeleteItems deletes multiple items from the database and writes response details to the http response writer
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
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// UpdateItem modifies a sinlge item in the database and writes the details to the http response writer, or returns an error if the item was not found
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

	product.Code = chi.URLParam(r, "code")

	response, err := app.UpdateDatabaseItem(product)
	if err != nil {
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// UpdateItems modifies multiple items in the database and writes new item details to the http response writer
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
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
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

// Search searches for an item by name in the database and writes all matching items to the response writer, or a not found error if there were no matches
func (app *Config) Search(w http.ResponseWriter, r *http.Request) {
	searchString := r.URL.Query().Get("item")
	searchString = strings.ToLower(searchString)

	response, err := app.SearchDatabaseItem(searchString)
	if err != nil {
		switch err.(type) {
		case *ErrNotFound:
			err := app.ErrorJSON(w, err, http.StatusNotFound)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		default:
			err := app.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				log.Printf("Problem with ErrorJSON: %s", err)
			}
		}
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "List Products",
		Data:    response,
	}

	err = app.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Problem with WriteJSON: %s", err)
	}
}
