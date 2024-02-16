package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type searchQuery string

var handlerTests = []struct {
	name               string
	url                string
	method             string
	product            any
	expectedStatusCode int
	searchTest         bool
}{
	{"Get Items", "/getitems", "GET", Product{}, http.StatusOK, false},
	{"Get Item Not Found", "/getitem", "GET", Product{}, http.StatusNotFound, false},
	{"Get Item", "/getitem", "GET", Product{"A12T-4GH7-QPL9-3N4M", "lettuce", 3.46}, http.StatusOK, false},
	{"Get Item 2", "/getitem", "GET", Product{"E5T6-9UI3-TH15-QR88", "Peach", 2.99}, http.StatusOK, false},

	{"Add Item Not Found", "/additem", "POST", Product{}, http.StatusBadRequest, false},
	{"Add Item", "/additem", "POST", Product{"A12T-4GH7-QPL9-9999", "Spinach", 1.29}, http.StatusOK, false},
	{"Update Item Not Found", "/updateitem", "POST", Product{}, http.StatusNotFound, false},
	{"Update Item", "/updateitem", "POST", Product{"A12T-4GH7-QPL9-9999", "Tuna", 1.29}, http.StatusOK, false},
	{"Delete Item Not Found", "/deleteitem", "POST", Product{}, http.StatusNotFound, false},
	{"Delete Item", "/deleteitem", "POST", Product{"A12T-4GH7-QPL9-9999", "Tuna", 1.29}, http.StatusOK, false},

	{"Add Items Not Found", "/additems", "POST", Products{}, http.StatusBadRequest, false},
	{"Add Items", "/additems", "POST", Products{
		[]Product{
			{
				Code:  "A12T-4GH7-QPL9-9997",
				Name:  "Cauliflower",
				Price: 100.99,
			},
			{Code: "A12T-4GH7-QPL9-9998",
				Name:  "Broccoli",
				Price: 1.33},
		},
	}, http.StatusOK, false},
	{"Update Items Not Found", "/updateitem", "POST", Products{}, http.StatusBadRequest, false},
	{"Update Items", "/updateitems", "POST", Products{
		[]Product{
			{
				Code:  "A12T-4GH7-QPL9-9997",
				Name:  "Twinkies",
				Price: 100.99,
			},
			{Code: "A12T-4GH7-QPL9-9998",
				Name:  "Slushies",
				Price: 1.33},
		},
	}, http.StatusOK, false},
	{"Delete Items Not Found", "/deleteitem", "POST", Products{}, http.StatusBadRequest, false},
	{"Delete Items", "/deleteitems", "POST", Products{
		[]Product{
			{
				Code:  "A12T-4GH7-QPL9-9997",
				Name:  "Twinkies",
				Price: 100.99,
			},
			{Code: "A12T-4GH7-QPL9-9998",
				Name:  "Slushies",
				Price: 1.33},
		},
	}, http.StatusOK, false},
	{"Search Items", "/search", "GET", "green", http.StatusOK, true},
}

func TestHandlers(t *testing.T) {
	var app Config
	err := app.init()
	if err != nil {
		log.Fatal(err)
	}
	routes := app.routes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _, e := range handlerTests {
		if e.searchTest == false {
			rb, err := json.Marshal(e.product)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			req, err := http.NewRequest(e.method, ts.URL+e.url, strings.NewReader(string(rb)))
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			resp, err := ts.Client().Get(ts.URL + e.url + "?item=green")
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
