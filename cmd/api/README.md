# Architecture

The project has been split out into several different files with specific purposes.

1. The server.go file calls the functions to register routes and middleware, intializes the config struct (including mutexes and database intialization), and starts the main web server.
2. The routes.go file registers the routes for the various API calls along with some basic middleware using chi and returns a server mux.
4. The handlers.go file gontains the handlers specified by routes.go. It uses the database and json helpers for the bulk of the work and returns json responses containing the response details and appropriate http status codes.
5. The helpers.go file contains JSON helpers to read and write JSON responses and errors.
6. The database.go file contains the structs and functions to work with the product database for the supermarket api.

## Server
The server.go file is the main web server. I tried to keep the server file as clean as possible and move various related functions into dedicated files for each service. 

The server has a Config struct to provide access to the in-memory database, mutexes, and various other settings. The purpose of the server.go file is to intialize this config struct along with the database, and then start the web server.

## Database
The database file contains the go structs for products/items and a dictionary of products that is used as the database for the API. 

### Product Struct
```
type Product struct {
	Code  string  `json:"code"`
	Name  string  `json:"name,omitempty"`
	Price float32 `json:"price,omitempty"`
}

```
The product struct contains three varaibles:
- A Code that is sixteen case-insensitive alphanumeric characters separated by dashes into four groups. The code format is verified using regex.
- A Name that is an alphnumeric string. The name is converted into lowercase before being used or stored in the dictionary for consistency. Names can also be used to search for items in the database using the /search call.
- A Price  of type float32. The price type has a stringer function so that it outputs two decimal places when printed.

There is a NewProduct function which takes a Code, Name, and Price as input, verifies the format, converts to upper/lower case as needed, and returns a correctly formatted Product. There are also VerifyProduct, VerifyCode, VerifyName, and VerifyPrice functions that will check the validity of each parameter and return useful errors.

There is also a Products struct that is simply a slice of Product structs.
```
type Products struct {
	Products []Product `json:"products"`
}
```

### Database Struct
The database struct is a map that uses the Product Code as the key so items can be looked up in constant time. The value for the dictionary is a Prodcut struct. 
```
Database = make(map[string]Product)
```

It is inialized with the products specified in the gist:
```
	initalProducts := []Product{
		{"a12T-4GH7-QPL9-3N4M", "Lettuce", 3.46}, //lower case a to test toUpper
		{"E5T6-9UI3-TH15-QR88", "Peach", 2.99},
		{"YRT6-72AS-K736-L4AR", "Green Pepper", 0.79},
		{"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", 3.59},
	}
```
#### Functions
The database uses a variety of functions to retrieve, add, update, and delete items. Importantly, the database functions have the Config struct as the receiver, which provides access to the mutexes and various other settings that can be configured. The first two lines for all database functions are Mutex.Lock() and defer Mutex.Unloc(), which ensures the database is safe and there are no race conditions.

##### Get, Add, Update, and Delete Item
The various functions for working with an individual Product all work in a similar way. They take in a Product and return a Product and error. 

In the case of Add and Update, the product that is passed into the function must contain a Code, Name, and Price. In the case of Get and Delete, only a code is required. When an operation succeeds, error is returned as nil. 

The general workflow for each function is to:
1. Lock the database mutex
2. Verify the product details that were passed in
3. Perform operations in the database
4. Check for errors 
5. Unlock the database mutex
6. Return the product details

There is a custom error type: "NewErrNotFound", which is used to identify errors caused by trying to get, update, or delete an item that is not in the database. This error type is checked by other services to determine if a database call failed due to a badly formed request or other issue, and return appropriate http status codes.

##### Get, Add, Update, and Delete Items
The various functions for working with mutliple items all take in a Products struct, which is a slice of Product structs, and return details about the Products and potentially an error. The execption to this is getitems, which does not take in any parameters since it returns all products in the database.

In the case of Add and Update, the products that are passed into the function must contain a Code, Name, and Price. In the case of Get and Delete, the slice of Product structs only require the code. When an operation succeeds, error is returned as nil. 

The general workflow for each function is to:
1. Lock the database mutex
2. Createa a response variable
3. Iterate over each product
   - Verify the product details that were passed in
   - Perform operations in the database
   - Check for errors 
   - Append the details to the response variable
7. Unlock the database mutex
8. Return the response details

There is a custom error type: "NewErrNotFound", which is used to identify errors caused by trying to get, update, or delete an item that is not in the database. This error type is checked by other services to determine if a database call failed due to a badly formed request or other issue, and return appropriate http status codes.

##### Search Items
The search function allows the database to be searched by name. The search string is converted to lower case, then each product in the database is checked to see if the name contains a match to the search string. Any matches are appended to a Products struct containing the response.

All matches are returned. If none are found, a not found status is returned.

## Routes
The routes file uses the chi router to register some basic middleware and then specifies handlers for various GET and POST requests such as getitem and updateitem. 

## Handlers
The handlers file contains all the functions registered by the routes file. Each handler takes a http response writer and returns an http request. The handlers make heavy use of the database functions and json helpers for the bulk of the work. 
- The handlers return a 404 Not Found status if an attempt is made to delete/update an item that doesn't exist, and a 400 Bad Request status codes if the JSON from the API call is not valid.
- Generally, the handlers return the product details to reflect any changes that were made. I was more verbose than usual with the JSON responses to make it easier to demo the API functionality.

## Helpers
The helpers are just functions that take in the http respone writer and rquest along with any data and attempt to marshall/unmarshall between JSON and Go structs. his abstraction keeps the handlers a lot cleaner. All of the helpers use the JSONResponse struct, which contains a bool value indicating whether there was an error, a short message, and the data containing the json response body.

- ReadJSON takes in a data variable of type any and attempts to unmarshal JSON from the http request body into the data variable. It uses a switch statement to determine details of what may have gone wrong if the unmarshall fails and return helpful error codes.
- WriteJSON takes in an http response writer, status code, data, and optional http headers. It marshalls the data into JSON and writes it to the http response along with the specified status and any provided headers.
- ErrorJSON is used to return error status codes when API calls fail or return errors. It takes an http response writer, error message, and status and writes the JSON to the response writer.


# Improvements
- Pagination for the items functions, such as getitems.
- Partial failiures on functions with multiple items.
- URLs should be something like /api/v1/
