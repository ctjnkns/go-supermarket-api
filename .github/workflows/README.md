
# Workflow Settings
Currently most workflows are configured to run when a pull request is created. Merges are still allowed in the case where workflows fail since this app is still in development. Commits directly to the main branch are not permittted.

# Workflows

## Lint Tests
The Lint test workflow runs the [golangci-lint](golangci-lint) tool using the configured linters and settings specified in the .golangic.toml file.

## Unit Tests
There is extensive unit testing using Go's built in testing library. Each file has an associated test file where each function is tested for expected inputs and ouptus and to confirm that the code is safe and consistent. 

For the main server file, the functions to initialize the server and database with default values are tested to make sure there will be no errors during setup:
```
func Test_Init(t *testing.T) {
	var app Config
	err := app.init()
	if err != nil {
		t.Errorf("Failed init")
	}
}
```

For http handlers, a variety of errors are tested to confirm the http status codes are working as expected:
```
	{"Get Items", "/getitems", "GET", Product{}, http.StatusOK, false},
	{"Get Item Bad Request", "/getitem", "GET", Product{}, http.StatusBadRequest, false},
	{"Get Item Not Found", "/getitem", "GET", Product{"E5T6-9UI3-TH15-QRZZ", "Peach", 2.99}, http.StatusNotFound, false},
```

For the json helper library, a variety of verbose errors are tested to provide as much useful information as possible.
```
{
	{name: "good json", json: `{"foo": "bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: false},
	{name: "badly formatted json", json: `{"foo":}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "incorrect type", json: `{"foo": 1}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "two json files", json: `{"foo": "1"}{"alpha": "beta"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "empty body", json: ``, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "syntax error in json", json: `{"foo": 1"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "unknown field in json", json: `{"fooo": "1"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "allow unknown fields in json", json: `{"fooo": "1"}`, errorExpected: false, maxSize: 1024, allowUnknown: true},
	{name: "missing field name", json: `{jack: "1"}`, errorExpected: true, maxSize: 1024, allowUnknown: true},
	{name: "file too large", json: `{"foo": "bar"}`, errorExpected: true, maxSize: 5, allowUnknown: true},
	{name: "file too large", json: `Hello, world!`, errorExpected: true, maxSize: 1024, allowUnknown: true},
}
```

Unit tests for the database still need to be implemented. 

## Build
The build workflow builds the docker container based on the docker-compose.yml and Dockerfile configratuion. A curl request is run to confirm that the server is responding and the database was initialized with default data as expected.
```yml
            - name: Build
              run: docker compose up -d
            - name: Test-API
              run: |
                curl --request GET 'http://localhost:8080/getitem' \
                    --header 'Content-Type: application/json' \
                    --data '{ "code": "A12T-4GH7-QPL9-3N4M" }'
```

## Publish
The publish build logs into dockerhub using secrets stored in the Github secrets repository, and then builds and publishes the docker image to dockerhub. This workflow runs manually but can be updated in the future to run on a commit to main or when a new module/release is tagged. It can also be configured to require a successful run of the build workflow to confirm that the dockerfile is correct. 
