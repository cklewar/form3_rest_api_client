# Author
Christian Klewar

# Description
Simple and concise REST API client library to interact with Form3 REST API written in Go and suitable for use in another software project. Current supported operations are:

* Create (POST)
* Fetch (GET)
* Delete (DELETE)

# Technical decisions
Contain documentation of your technical decisions.

# Requirements

```go
golang.org/x/mod v0.4.2 // indirect
golang.org/x/net v0.0.0-20210315170653-34ac3e1c2000 // indirect
golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
golang.org/x/sys v0.0.0-20210315160823-c6e025ad8005 // indirect
golang.org/x/tools v0.1.0 // indirect
```

# Implementation

## Operations
Operations are defined in __Operations__ interface. __Operations__ type define common method sets. Type __Parameters__ therefore, is said to implement __Operations__ interface by implementing its methods. In that light, __Operations__ interface enable us compose custom types that have a common behavior.

```go
type Operations interface {
	Create(id string) (int, []byte)
	Delete(id string, version int) int
	Fetch(id string) (int, []byte)
}
```

## Parameters 
Parameters are implemented as struct. Parameters struct defines mandatory fields to be used to generate URI for later use. Parameters struct and it's fields are kept public except __uri__ field which needs to be generated out of other field data. __uri__ is the final string used to connect to REST endpoint. 

```go
type Parameters struct {
	Timeout     time.Duration // HTTP wait timeout
	Protocol    string        // HTTP or HTTPS
	Host        string        // IP or DNS name of target host
	Port        string        // TCP port number on target host
	BaseURI     string        // base URI e.g. "/v1/organisation/" or "/v1/transaction/payments"
	ContentType string        // Header content type e.g. application/vnd.api+json
	Resource    string        // API resource endpoint e.g. account, claims
	uri         string        // not public needs to be generated
}
```

### Fields
- Timeout defines HTTP wait timeout for ```http.Client{}```
- Protocol defines protocol schema e.g. http or https        
-	Host defines target IP or DNS name. This is the API server address
-	Port defines targets TCP port number. This is the API service listening port
-	BaseURI defines base URI e.g. __/v1/organisation/__ or __/v1/transaction/payments__
-	ContentType defines request header content type e.g. __application/vnd.api+json__
-	Resource defines API resource endpoint e.g. __account__, __claims__
-	uri defines final uri connection string. Intentionally not public needs to be generated

## API Client
API Client is implemented as struct embedding unnamed __Parameters__ struct. API client can be created using __NewClient()__ function. 

```go
// Client is a struct which embeds Parameters
type Client struct {
	Parameters
}
```

Client struct implements __Operations__ interface.

## Client API Constructor
To create a REST API client one can use __NewClient()__ function. Client constructor will take __Parameters__ struct to initialize API client. 

```go
func NewClient(o Parameters) (bool, *Client) {}
```

__NewClient__ function returns status boolean and initialized __Client__ struct. 

## Operation Methods

### Create
Create new resource with __input__ data on REST endpoint.

```go
func (c *Client) Create(input []byte) (int, []byte) {}
```
#### Return values
__int__: Integer < 0 if error occured or HTTP status code e.g. 404, 201  
__[]byte__: API server response body


### Delete
Delete resource with __id__ and __version__.
```go
func (c *Client) Delete(id string, version int) int {}
```

#### Return value
__int__: Integer < 0 if error occured or HTTP status code e.g. 404, 201 

### Fetch 

```go
func (c *Client) Fetch(id string) (int, []byte) {}
```

#### Return values
__int__: Integer < 0 if error occured or HTTP status code e.g. 404, 201  
__[]byte__: API server response body


# Usage

## Import client library

```go
"github.com/cklewar/form3_clib/api/client"
```

## Initialize parameters

```go
// Initialize API client parameters
parameters := api.Parameters{
	Host:     "192.168.2.50",
	BaseURI:  "/v1/organisation/",
	Resource: "accounts",
}
```

## Create API client
```go
// Construct API Client
status, c := client.NewClient(parameters)
```

## Operations

### Create
In this example we use input file as JSON source for create operation.

```go
cwd, err := os.Getwd()
path := filepath.Join(cwd, "data/org_acc_create.json")
createInputData, err := ioutil.ReadFile(path)

if err != nil {
    log.Fatal(err)
}
```

Create a resource by calling create method on client.

```go
status, createResp := c.Create(createInputData)
fmt.Println("Status: ", status)
fmt.Println(client.JSONPrettyPrint(createResp))
```

Output

```bash
Status:  201
{
  "data": {
    "attributes": {
      "alternative_bank_account_names": null,
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "base_currency": "GBP",
      "bic": "NWBKGB22",
      "country": "GB"
    },
    "created_on": "2021-03-16T21:43:15.798Z",
    "id": "de37f385-1604-7c5b-a1e5-3173ea9cc2db",
    "modified_on": "2021-03-16T21:43:15.798Z",
    "organisation_id": "ec0bd2f5-d6f7-54b3-b687-acd21cdde71b",
    "type": "accounts",
    "version": 0
  },
  "links": {
    "self": "/v1/organisation/accounts/de37f385-1604-7c5b-a1e5-3173ea9cc2db"
  }
}
```

### Fetch
Fetch a resource by it's __id__. 

```go
status, fetchResp := c.Fetch("ad27e265-9604-4b4b-a0e5-3003ea9cc4db")
fmt.Println("Status: ", status)
fmt.Println(client.JSONPrettyPrint(fetchResp))
```

Output

```bash
Status: 200
{
  "data": {
    "attributes": {
      "alternative_bank_account_names": null,
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "base_currency": "GBP",
      "bic": "NWBKGB22",
      "country": "GB"
    },
    "created_on": "2021-03-15T15:49:57.788Z",
    "id": "ad27e265-9604-4b4b-a0e5-3003ea9cc4db",
    "modified_on": "2021-03-15T15:49:57.788Z",
    "organisation_id": "eb0bd3f5-c3f5-44b2-b677-acd21cdde71b",
    "type": "accounts",
    "version": 0
  },
  "links": {
    "self": "/v1/organisation/accounts/ad27e265-9604-4b4b-a0e5-3003ea9cc4db"
  }
}
```

### Delete
Delete a resource by it's __id__ and __version number__.

```go
status = c.Delete("UUID", 0)
fmt.Println("Status: ", status)
```

Output
```bash
Status:  204
```

## Response
Client library provides structured response data within response package. Response package defines e.g. __Organisation Account__ data structures for unmarshaling API server JSON response data into structs. 

### Example

```go
var account response.OrganisationAccountData
json.Unmarshal([]byte(body), &account)
fmt.Println("Id:", account.Data.ID)
fmt.Println("Type:", account.Data.Type)
```

Output

```bash
Id: ad27e265-9604-4b4b-a0e5-3003ea9cc4db
Type: accounts
```

# Testing
Be well tested to the level you would expect in a commercial environment.
Have tests that run from docker-compose up - our reviewers will run docker-compose up to assess if your tests pass.
## Linter
## Unit tests
## Functional tests
## Smoke tests
