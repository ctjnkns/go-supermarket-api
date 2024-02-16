package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/getitem", app.GetItem)
	mux.Get("/getitems", app.GetItems)

	mux.Post("/additem", app.AddItem)
	mux.Post("/additems", app.AddItems)

	mux.Post("/deleteitem", app.DeleteItem)
	mux.Post("/deleteitems", app.DeleteItems)

	mux.Post("/updateitem", app.UpdateItem)
	mux.Post("/updateitems", app.UpdateItems)

	mux.Get("/search", app.Search)

	return mux
}
