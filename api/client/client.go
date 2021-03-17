package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Return error state value
const error = -1

// Version is the current library's version: sent with User-Agent
const Version = "0.1"

// Parameters default values
const (
	defaultProtocol    = "http"
	defaultTimeout     = time.Second * 10
	defaultPort        = "8080"
	defaultContentType = "application/vnd.api+json"
)

// Operations defined in public interface
type Operations interface {
	Create(id string) (int, []byte)
	Delete(id string, version int) int
	Fetch(id string) (int, []byte)
}

// Parameters public structure
type Parameters struct {
	Timeout     time.Duration // HTTP wait timeout. Default is time.Second * 10
	BaseURI     string        // base URI e.g. "/v1/organisation/", "/v1/transaction/". Need trailing slash!. Mandatrory field
	ContentType string        // Header content type. Default is application/vnd.api+json
	Resource    string        // API resource endpoint e.g. account, claims. Mandatory field
}

// Client is a struct which embeds Parameters
type Client struct {
	Parameters
	protocol string // HTTP or HTTPS. Default is HTTP
	host     string // IP or DNS name of target host. Mandatory field
	port     string // TCP port number on target host. Default is 8080
	uri      string // not public needs to be generated
}

// Create (POST) a new ressource to <request.uri>.
// Return tuple(<StatusCode>, <response Body>).
// In case of error return <error> as status.
func (c *Client) Create(input []byte) (int, []byte) {
	var body []byte

	fmt.Printf("Creating new resource at URI <%s>\n", c.uri)
	req, err := http.NewRequest(http.MethodPost, c.uri, bytes.NewBuffer(input))

	if err != nil {
		log.Fatal("Error reading request. ", err)
		return error, body
	}

	req.Header.Add("Content-type", c.ContentType)
	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error reading response. ", err)
		return error, body
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error reading body. ", err)
		return error, body
	}

	return resp.StatusCode, body
}

// Delete a resource by <id> and return StatusCode.
// <StatusCode> code in case of error < 0
func (c *Client) Delete(id string, version int) int {
	res := c.uri + id + "?version=" + strconv.Itoa(version)
	fmt.Printf("Deleting resource at URI <%s>\n", res)
	req, err := http.NewRequest(http.MethodDelete, res, nil)

	if err != nil {
		log.Fatal("Error reading request. ", err)
		return error
	}

	req.Header.Add("Content-type", c.ContentType)
	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error reading response. ", err)
		return error
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

// Fetch (GET) data from API. Parameter is a resource's <id>.
// Return tuple(<StatusCode>, <response Body>).
// In case of error return <error> as status.
func (c *Client) Fetch(id string) (int, []byte) {
	var body []byte
	res := c.uri + id
	fmt.Printf("Fetching from URI <%s>\n", res)
	req, err := http.NewRequest(http.MethodGet, res, nil)

	if err != nil {
		log.Fatal("Error reading request. ", err)
		return error, body
	}

	req.Header.Add("Content-type", c.ContentType)
	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error reading response. ", err)
		return error, body
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error reading body. ", err)
		return error, body
	}

	return resp.StatusCode, body
}

// GetObjID is a support method extracting <id> out of JSON data
// GetObjID input parameter is byte array data from API response
// Method handels byte array data as unstructured data by leveraging generic interface map
// Method returns id as string value. In case of error empty string returned
func (c *Client) GetObjID(input []byte) string {
	var id string
	// Declared array map of string with empty interface (unstructured input data)
	// which will hold the value of the parsed json. Parsing embedded object <id> in JSON and return.
	var raw map[string]interface{}
	err := json.Unmarshal(input, &raw)

	if err != nil {
		log.Fatal(err)
		return id
	}

	raw = raw["data"].(map[string]interface{})
	id = raw["id"].(string)

	return id
}

// GetObjVersion is a support method extracting <version> out of JSON data
// GetObjVersion input parameter is byte array data from API response
// Method handels byte array data as unstructured data by leveraging generic interface map
// Method returns version as integer value. In case of error <error> returned
func (c *Client) GetObjVersion(input []byte) int {
	// Declared array map of string with empty interface (unstructured input data)
	// which will hold the value of the parsed json. Parse embedded object <version> in JSON and return.
	var raw map[string]interface{}
	err := json.Unmarshal(input, &raw)

	if err != nil {
		log.Fatal(err)
		return error
	}

	raw = raw["data"].(map[string]interface{})
	// Converting version type from float to integer.
	// Investigate why is value mapped as float when using empty interface?
	version := int(raw["version"].(float64))

	return version
}

//
// Default value methods start
//
func (p *Parameters) contentTypeBase() string {
	base := p.ContentType
	if p.ContentType == "" {
		base = defaultContentType
	}
	return base
}

func (p *Parameters) timeoutBase() time.Duration {
	base := p.Timeout
	if p.Timeout == 0 {
		base = defaultTimeout
	}
	return base
}

func (c *Client) protocolBase(protocol string) string {
	if protocol == "" {
		return defaultProtocol
	}
	return protocol
}

func (c *Client) portBase(port string) string {
	if port == "" {
		return defaultPort
	}
	return port
}

// NewClient constructor with default values check
func NewClient(host string, port string, protocol string, p Parameters) (bool, *Client) {
	var c Client

	if host == "" {
		log.Fatal("Host parameter not set")
		return false, &c
	}

	c = Client{
		host:       host,
		port:       c.portBase(port),
		protocol:   c.protocolBase(protocol),
		Parameters: p,
	}

	if c.BaseURI == "" {
		log.Fatal("BaseURI parameter not set")
		return false, &c
	}

	if c.Resource == "" {
		log.Fatal("Resource parameter not set")
		return false, &c
	}
	// Check for setting parameters default value
	c.Timeout = p.timeoutBase()
	c.ContentType = p.contentTypeBase()

	//Set final uri connect string
	c.uri = c.protocol + "://" + c.host + ":" + c.port + c.BaseURI + c.Resource + "/"

	return true, &c
}

// JSONPrettyPrint prints JSON data in a more readable way on terminal.
func JSONPrettyPrint(body []byte) string {
	dst := &bytes.Buffer{}

	if err := json.Indent(dst, body, "", "  "); err != nil {
		log.Fatal(err)
	}
	return dst.String()
}
