# go-supermarket-api

## Setup
- Run `docker compose up -d` to run the docker container.
- Run `docker compuse build` to re-build the image.
- The postman collection can be imported to easily test.

## Usage

### Working with Single Items

#### Getting a specific item

A specific item can be looked up by it's product code.
```shell
curl --request GET 'http://localhost:8080/getitem' \
--header 'Content-Type: application/json' \
--data '{
        "code": "A12T-4GH7-QPL9-9999"
}'
```

The details of the itme are returned along with status code 200.
```json
{
    "error": false,
    "message": "List product: A12T-4GH7-QPL9-9999",
    "data": {
        "code": "A12T-4GH7-QPL9-9999",
        "name": "spinach",
        "price": 1.29
    }
}
```

If the item is not found, a 404 not found error is returned.
```json
{
    "error": true,
    "message": "Item not found: A12T-4GH7-QPL9-9999"
}
```

#### Adding an item

An item can be created by sending the item details in a json post request.
```shell
curl 'http://localhost:8080/additem' \
--header 'Content-Type: application/json' \
--data '{
        "code": "A12T-4GH7-QPL9-9999",
        "name": "Spinach",
        "price": 1.29
}'
```

The details of the itme are returned along with status code 200. 
```json
{
    "error": false,
    "message": "Added A12T-4GH7-QPL9-9999",
    "data": {
        "code": "A12T-4GH7-QPL9-9999",
        "name": "spinach",
        "price": 1.29
    }
}
```
If there is an error in the format of the json or the item details, an error message is returned along with status code 400.

```json
{
    "error": true,
    "message": "Unable to create product: Unable to verify code: Incorrect Product Code Format"
}
```

#### Updating an existing item
```shell
curl 'http://localhost:8080/updateitem' \
--header 'Content-Type: application/json' \
--data '{
        "code": "A12T-4GH7-QPL9-9999",
        "name": "Tuna",
        "price": 1.29
}'
```

A short message and the details of the updated item are returned.
```json
{
    "error": false,
    "message": "Updated A12T-4GH7-QPL9-9999",
    "data": {
        "code": "A12T-4GH7-QPL9-9999",
        "name": "tuna",
        "price": 1.29
    }
}
```

If the request is not properly formatted or the item is not found, a 400 or 404 status code is returned along with an error message.

```shell
{
        "code": "A12T-4GH7-QPL9-9997",
        "name": "Tuna",
        "price": 1.29
}
```

```json
{
    "error": true,
    "message": "Item not found: A12T-4GH7-QPL9-9997"
}
```

#### Deleting an item
```shell
curl 'http://localhost:8080/deleteitem' \
--header 'Content-Type: application/json' \
--data '{
        "code": "A12T-4GH7-QPL9-9999"
}'
```

A short message, success status code (200), and the details of the deleted item are returned.
```json
{
    "error": false,
    "message": "Deleted product: A12T-4GH7-QPL9-9999",
    "data": {
        "code": "A12T-4GH7-QPL9-9999",
        "name": "tuna",
        "price": 1.29
    }
}
```

If the item is not found, a 404 is returned along with an error message
```json
{
    "error": true,
    "message": "Item not found: A12T-4GH7-QPL9-9997"
}
```

### Working with Mutliple Items

#### Getting all items in the dtatabse
```shell
curl  'http://localhost:8080/getitems'
```

The response will contain a bool value indicating whehter or not there was an error, a short message, and a list of products in JSON. A successfull request returns status code 200.
```json
{
    "error": false,
    "message": "Listing products",
    "data": {
        "A12T-4GH7-QPL9-3N4M": {
            "code": "A12T-4GH7-QPL9-3N4M",
            "name": "lettuce",
            "price": 3.46
        },
        "E5T6-9UI3-TH15-QR88": {
            "code": "E5T6-9UI3-TH15-QR88",
            "name": "peach",
            "price": 2.99
        },
        "TQ4C-VV6T-75ZX-1RMR": {
            "code": "TQ4C-VV6T-75ZX-1RMR",
            "name": "gala apple",
            "price": 3.59
        },
        "YRT6-72AS-K736-L4AR": {
            "code": "YRT6-72AS-K736-L4AR",
            "name": "green pepper",
            "price": 0.79
        }
    }
}
```

#### Adding multiple items
```shell
curl 'http://localhost:8080/additems' \
--data '{
    "products": [
        {
                "code": "A12T-4GH7-QPL9-9997",
                "name": "Cauliflower",
                "price": 100.99
        },
        {
                "code": "A12T-4GH7-QPL9-9998",
                "name": "Broccoli",
                "price": 1.33
        }
    ]
}'
```

Status code 200 along with some details are returned.
```json
{
    "error": false,
    "message": "Added multiple products",
    "data": {
        "products": [
            {
                "code": "A12T-4GH7-QPL9-9997",
                "name": "cauliflower",
                "price": 100.99
            },
            {
                "code": "A12T-4GH7-QPL9-9998",
                "name": "broccoli",
                "price": 1.33
            }
        ]
    }
}
```

If the request is badly formed an error is returned.
```json
{
    "error": true,
    "message": "Unable to create product: Unable to verify code: Incorrect Product Code Format"
}
```

#### Updating existing items

```shell
curl 'http://localhost:8080/updateitems' \
--header 'Content-Type: application/json' \
--data '{
    "products": [
        {
                "code": "A12T-4GH7-QPL9-9997",
                "name": "Twinkies",
                "price": 100.99
        },
        {
                "code": "A12T-4GH7-QPL9-9998",
                "name": "Slushies",
                "price": 1.33
        }
    ]
}'
```

A success status code and item details are returned.
```json
{
    "products": [
        {
                "code": "A12T-4GH7-QPL9-9997",
                "name": "Twinkies",
                "price": 100.99
        },
        {
                "code": "A12T-4GH7-QPL9-9998",
                "name": "Slushies",
                "price": 1.33
        }
    ]
}
```

If the item does not exist, an error is returned.
```json
{
    "error": true,
    "message": "item not found: A12T-4GH7-QPL9-8997"
}
```
#### Deleting items
Multiple items can be deleted at a time using their product codes.
```shell
curl 'http://localhost:8080/deleteitems' \
--header 'Content-Type: application/json' \
--data '{
    "products": [
        {
                "code": "A12T-4GH7-QPL9-9997"
        },
        {
                "code": "A12T-4GH7-QPL9-9998"
        }
    ]
}'
```

Details of the deleted items are returned.
```json
{
    "error": false,
    "message": "Deleted multiple products",
    "data": {
        "products": [
            {
                "code": "A12T-4GH7-QPL9-9997",
                "name": "Twinkies",
                "price": 100.99
            },
            {
                "code": "A12T-4GH7-QPL9-9998",
                "name": "Slushies",
                "price": 1.33
            }
        ]
    }
}
```

If any of the items don't exist, an error is returned.
```json
{
    "error": true,
    "message": "Item not found: "
}
```
### Searching

Items can be searched for using a query parameter. Partial matches will result in a success, and a list of all matches is returned.
```shell
curl 'http://localhost:8080/search?item=g'
```

```json
{
    "error": false,
    "message": "List Products",
    "data": {
        "products": [
            {
                "code": "YRT6-72AS-K736-L4AR",
                "name": "green pepper",
                "price": 0.79
            },
            {
                "code": "TQ4C-VV6T-75ZX-1RMR",
                "name": "gala apple",
                "price": 3.59
            }
        ]
    }
}
```


```shell
curl 'http://localhost:8080/search?item=green'
```

```json
{
    "error": false,
    "message": "List Products",
    "data": {
        "products": [
            {
                "code": "YRT6-72AS-K736-L4AR",
                "name": "green pepper",
                "price": 0.79
            }
        ]
    }
}
```

```shell
curl 'http://localhost:8080/search?item=greenapples'
```

```json
{
    "error": true,
    "message": "item not found"
}
```
