- [Author](#author)
- [Description](#description)
- [Technical decisions](#technical-decisions)
- [Requirements](#requirements)
- [Implementation](#implementation)
  - [APIInterface](#apiinterface)
  - [Parameters](#parameters)
    - [Fields](#fields)
  - [Client](#client)
    - [Fields](#fields-1)
  - [Response](#response)
    - [Fields](#fields-2)
  - [Client constructor](#client-constructor)
  - [Operation Methods](#operation-methods)
    - [Create](#create)
      - [Return values](#return-values)
    - [Delete](#delete)
      - [Return value](#return-value)
    - [Fetch](#fetch)
      - [Return values](#return-values-1)
- [Usage](#usage)
  - [Docker Compose](#docker-compose)
  - [Programatic](#programatic)
  - [Import client library](#import-client-library)
  - [Initialize parameters](#initialize-parameters)
  - [Create client](#create-client)
  - [Operations](#operations)
    - [Create](#create-1)
    - [Fetch](#fetch-1)
    - [Delete](#delete-1)
    - [Parameters](#parameters-1)
  - [Response](#response-1)
    - [Example](#example)
- [Testing](#testing)
  - [Unit](#unit)
  - [Linter and Formatter](#linter-and-formatter)
  - [Functional](#functional)
  - [Integration](#integration)

# Author
Christian Klewar

Go language experience is entry level / novice

# Description
Simple and concise REST API client library to interact with Form3 REST API written in Go and suitable for use in another software project. Current supported operations are:

* Create (POST)
* Fetch (GET)
* Delete (DELETE)

# Technical decisions
This chapter describes the technical decisions which have been made regarding client API design and integration.
Given the task client API should be

* simple and concise
* a client library suitable for use in another software project

following higher level technical descisions have been made:

* Abstraction
  * Hide complexity by introducing abstraction. To achive this client library leverages GO language specific constructs like method receivers, structs
* Usage
  * Client library should be as easy as possible to use. To achive this client library introduces client initilizer function
  * The initilaizer function has been build in a way to support reusable client "instance" 
  * User will define API server specific parameters like IP, port and protocol schema only once through entire API server usage
  * User defines more `variable` parameters within __parameters__ struct. `variable` parameters for example are a __resource__ or a __base uri__  which can be changed any time through entire API server usage
  * The library makes usage of default values whenever it makes sense and possible. 
* Extendability
  * Client library should be extendabale. To achive this client library introduces proper API interface which makes it possible to:
      * add new functionality
      * add missing functionality
      * introduce backwards compatability 
      * introduce versioning support
* Simplicity
  * Client library should be simple and consice. To achive this client library 
    * implementes proper functions and methods to abstract complexity
    * provides clean return value structure 
    * provides the possibility to build structured data out of given response data
* Standard library only
  * no 3rd party libraries were used 

# Requirements

```go
golang.org/x/mod v0.4.2 // indirect
golang.org/x/net v0.0.0-20210315170653-34ac3e1c2000 // indirect
golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
golang.org/x/sys v0.0.0-20210315160823-c6e025ad8005 // indirect
golang.org/x/tools v0.1.0 // indirect
```

# Implementation

## APIInterface
Common operations regarding Form3 REST API are defined in __APIInterface__ interface. Type __Client__ therefore, is said to implement __APIInterface__ interface by implementing __APIInterface__ methods. In that light, __APIInterface__ interface enables us to compose custom types that have a common behavior.

```go
// APIInterface is a public interface
type APIInterface interface {
	Create(input []byte, timeout time.Duration) (Response, error)
	Delete(id string, version int, timeout time.Duration) (Response, error)
	Fetch(id string, timeout time.Duration) (Response, error)
}
```

## Parameters 
Parameters are implemented as struct. Parameters struct defines mandatory fields to be used to generate URI for later use. Parameters struct kept private. Changing 
field values must be done through update methods. 

```go
type parameters struct {
	BaseURI     string // base URI e.g. "/v1/organisation/", "/v1/transaction/". Need trailing slash!. Mandatrory field
	ContentType string // Header content type. Default is application/vnd.api+json
	Resource    string // API resource endpoint e.g. account, claims. Mandatory field
}
```

### Fields
-	__BaseURI__ defines base URI e.g. __/v1/organisation/__ or __/v1/transaction/payments__
-	__ContentType__ defines request header content type e.g. __application/vnd.api+json__
-	__Resource__ defines API resource endpoint e.g. __account__, __claims__

__BaseURI__ and __Resource__ are mandatory fields and do not provide default values.

## Client
Client API is implemented as struct embedding unnamed __Parameters__ struct. API client can be created using __NewClient()__ function. 

```go
// Client is a struct which embeds parameters struct (unnamed)
type Client struct {
	parameters
	protocol string // HTTP or HTTPS. Default is HTTP
	host     string // IP or DNS name of target host. Mandatory field
	port     string // TCP port number on target host. Default is 8080
}
```

### Fields
- __protocol__ defines protocol schema e.g. http or https        
-	__host__ defines target IP or DNS name. This is the API server address
-	__port__ defines targets TCP port number. This is the API service listening port

Client struct implements __APIInterface__ interface.

```go
var _ APIInterface = (*Client)(nil) // Verify that *Client implements APIInterface
```

## Response
Response is used to return API server body data and HTTP response code

```go
// Response is used to return API server body data and according http response code
type Response struct {
	Body []byte
	Code int
}
```

### Fields
- __Body__ stores API server repsonse body data e.g. json data
- __Code__ stores API server repsonse code        


## Client constructor
To create a new client one can use __NewClient()__ function. 

* Function will take host, port and protocol
  * host is a mandatory field
  * port and protocol fields can be initialized with default values
* Function will take __Parameters__ struct to initialize "variable" data. 

```go
func NewClient(host string, port string, protocol string, p Parameters) (*Client, error){}
```

__NewClient__ function returns initialized __Client__ and __error__. 
__Client__ must implement __APIInterface__.

## Operation Methods

### Create
Create new resource with __input__ data and timeout set.

```go
Create(input []byte, timeout time.Duration) (Response, error){}
```
#### Return values
- __Response__: Response struct
- __error__: error


### Delete
Delete resource with __id__, __version__ and timeout set.
```go
Delete(id string, version int, timeout time.Duration) (Response, error) {}
```

#### Return value
- __Response__: Response struct
- __error__: error

### Fetch 
Fetch resource with __id__ and timeout set.
```go
Fetch(id string, timeout time.Duration) (Response, error) {}
```

#### Return values
- __Response__: Response struct
- __error__: error

# Usage

## Docker Compose
Run __docker-compose up__ from terminal. Entries starting with `client_1` are related to client api service.

Output:

```bash
client_1      | ..=== RUN   TestNewClient
client_1      | --- PASS: TestNewClient (0.00s)
client_1      | === RUN   TestCreate
client_1      | PATH: ../../examples/json/org_acc_create.json
client_1      | Creating new resource at URI <http://accountapi:8080/v1/organisation/accounts/>
client_1      | Deleting resource at URI <http://accountapi:8080/v1/organisation/accounts/de37f789-1604-7c5b-a1e5-3673ea9cc2db?version=0>
client_1      | --- PASS: TestCreate (0.02s)
client_1      | === RUN   TestFetch
client_1      | Creating new resource at URI <http://accountapi:8080/v1/organisation/accounts/>
client_1      | Fetching from URI <http://accountapi:8080/v1/organisation/accounts/de37f789-1604-7c5b-a1e5-3673ea9cc2db>
client_1      | Deleting resource at URI <http://accountapi:8080/v1/organisation/accounts/de37f789-1604-7c5b-a1e5-3673ea9cc2db?version=0>
client_1      | --- PASS: TestFetch (0.01s)
client_1      | === RUN   TestDelete
client_1      | Creating new resource at URI <http://accountapi:8080/v1/organisation/accounts/>
client_1      | Deleting resource at URI <http://accountapi:8080/v1/organisation/accounts/de37f789-1604-7c5b-a1e5-3673ea9cc2db?version=0>
client_1      | --- PASS: TestDelete (0.00s)
client_1      | === RUN   TestContentTypeBase
client_1      | --- PASS: TestContentTypeBase (0.00s)
client_1      | === RUN   TestTimeoutBase
client_1      | --- PASS: TestTimeoutBase (0.00s)
client_1      | === RUN   TestPortBase
client_1      | --- PASS: TestPortBase (0.00s)
client_1      | === RUN   TestProtocolBase
client_1      | --- PASS: TestProtocolBase (0.00s)
client_1      | PASS
client_1      | coverage: 78.0% of statements
client_1      | ok      github.com/cklewar/form3_rest_api_client/api/client     0.042s  coverage: 78.0% of statements
client_1      | ?       github.com/cklewar/form3_rest_api_client/api/response   [no test files]
client_1      | ?       github.com/cklewar/form3_rest_api_client/examples       [no test files]
```

## Programatic

## Import client library

```go
"github.com/cklewar/form3_clib/api/client"
```

## Initialize parameters

```go
// Initialize API client parameters
parameters := api.Parameters{
	BaseURI:  "/v1/organisation/",
	Resource: "accounts",
}
```

## Create client
Creating / initializing a new client is done by calling __NewClient()__ function. 
```go
// Construct API Client
c, err := client.NewClient("192.168.2.50", "", "", parameters)
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
createResp, err := c.Create(createInputData, 0)
fmt.Println("Error: ", err)
fmt.Println("ResponseCode: ", createResp.Code)
data, err := client.JSONPrettyPrint(createResp.Body)
fmt.Println(data)
```

Output

```bash
Error: nil
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
fetchResp, err := c.Fetch(id, 0)
fmt.Println("Error: ", err)
fmt.Println("ResponseCode: ", fetchResp.Code)
data, err = client.JSONPrettyPrint(fetchResp.Body)
fmt.Println(data)
```

Output

```bash
Error: nil
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
deleteResp, err := c.Delete(id, version, 0)
fmt.Println("Error: ", err)
fmt.Println("ResponseCode: ", deleteResp.Code)
```

Output
```bash
Error:  <nil>
ResponseCode:  204
```

### Parameters
Changing parameters by creating new parameters struct and calling __updateParameters()__ function. 

```go
parameters = client.Parameters{
  BaseURI:  "/v2/organisation/",
  Resource: "payments",
}

err := c.UpdateParameters(parameters)
```

Changing clients parameter values by accessing and assigning new values on client struct.

```go
c.UpdateContentType("NEW CONTENT TYPE")
c.UpdateResource("NEW RESOURCE")
c.UpdateBaseURI("NEW BASEURI")
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
* Have tests that run from docker-compose up
* Reviewers will run docker-compose up to assess if your tests pass

Desicion has been made to use `golang:latest` debian based for testing instead of `golang:alpine` which has smaller image size but issues have been faced with missing libs regarding unit testing.

## Unit
Unit test coverage is not complete. Approach is to catch as many use cases / corner cases as possible. At least every function / method should be unit tested.

Run unit tests with: 

```bash
go test -v ./api/client/ -coverprofile cover.out
```

## Linter and Formatter
* __golint__ used as linter and analyzes source code to flag programming errors, bugs, stylistic errors, and suspicious constructs.
* __gofmt__ used to shape the source code 

## Functional
Functional test would expect to get a specific value from the API server as defined by the requirements. Functional tests are not part of this repository. 

## Integration
Integration testing is done using micro service approach. Client API testing the interaction with the API server serivice and making sure that microservices work together as expected.
