package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Version is the current library's version: sent with User-Agent
const Version = "0.1"

// Default values
const (
	defaultProtocol    = "http"
	defaultTimeout     = time.Second * 10
	defaultPort        = "8080"
	defaultContentType = "application/vnd.api+json"
)

// ErrParamNotSet is used when mandatory parameter not set
var ErrParamNotSet = errors.New("parameter not set")

// APIInterface is a public interface
type APIInterface interface {
	Create(input []byte, timeout time.Duration) (Response, error)
	Delete(id string, version int, timeout time.Duration) (Response, error)
	Fetch(id string, timeout time.Duration) (Response, error)
}

// Parameters struct
type parameters struct {
	BaseURI     string // base URI e.g. "/v1/organisation/", "/v1/transaction/". Need trailing slash!. Mandatrory field
	ContentType string // Header content type. Default is application/vnd.api+json
	Resource    string // API resource endpoint e.g. account, claims. Mandatory field
}

// Client is a struct which embeds parameters struct
type Client struct {
	parameters
	protocol string // HTTP or HTTPS. Default is HTTP
	host     string // IP or DNS name of target host. Mandatory field
	port     string // TCP port number on target host. Default is 8080
}

var _ APIInterface = (*Client)(nil) // Verify that *Client implements APIInterface

// Response is used to return API server body data and according http response code
type Response struct {
	Body []byte
	Code int
}

func (c *Client) updateParameters(p parameters) error {

	err := c.updateBaseURI(p.BaseURI)

	if err != nil {
		return err
	}

	err = c.updateResource(p.Resource)

	if err != nil {
		return err
	}

	if p.ContentType == "" {
		c.ContentType = defaultContentType
	} else {
		c.ContentType = p.ContentType
	}

	return nil
}

func (c *Client) updateBaseURI(baseURI string) error {
	if baseURI == "" {
		return fmt.Errorf("%q: %w", "BaseURI", ErrParamNotSet)
	}

	c.BaseURI = baseURI
	return nil
}

func (c *Client) updateResource(resource string) error {
	if resource == "" {
		return fmt.Errorf("%q: %w", "BaseURI", ErrParamNotSet)
	}

	c.Resource = resource
	return nil
}

func (c *Client) updateURI() (string, error) {
	if c.BaseURI != "" && c.Resource != "" {
		//Set final uri connect string
		uri := c.protocol + "://" + c.host + ":" + c.port + c.BaseURI + c.Resource + "/"
		return uri, nil
	}
	return "", errors.New("Missing mandatory field")
}

// Create (POST) a new ressource to <request.uri>.
// Return tuple(<StatusCode>, <response Body>).
// Return error
func (c *Client) Create(input []byte, timeout time.Duration) (Response, error) {
	response := Response{
		Body: nil,
		Code: -1,
	}

	uri, err := c.updateURI()

	if err != nil {
		return response, err
	}

	fmt.Printf("Creating new resource at URI <%s>\n", uri)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(input))

	if err != nil {
		return response, err
	}

	req.Header.Add("Content-type", c.ContentType)
	client := &http.Client{Timeout: checkTimeout(timeout)}
	resp, err := client.Do(req)

	if err != nil {
		response.Code = resp.StatusCode
		return response, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}

	response.Body = body
	response.Code = resp.StatusCode

	return response, nil
}

// Delete a resource by <id> and return StatusCode.
// Return error
func (c *Client) Delete(id string, version int, timeout time.Duration) (Response, error) {
	response := Response{
		Body: nil,
		Code: -1,
	}

	uri, err := c.updateURI()

	if err != nil {
		return response, err
	}

	res := uri + id + "?version=" + strconv.Itoa(version)
	fmt.Printf("Deleting resource at URI <%s>\n", res)
	req, err := http.NewRequest(http.MethodDelete, res, nil)

	if err != nil {
		return response, err
	}

	req.Header.Add("Content-type", c.ContentType)
	client := &http.Client{Timeout: checkTimeout(timeout)}
	resp, err := client.Do(req)

	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	response.Code = resp.StatusCode

	return response, nil
}

// Fetch (GET) data from API. Parameter is a resource's <id>.
// Return Response struct
// Return error
func (c *Client) Fetch(id string, timeout time.Duration) (Response, error) {
	response := Response{
		Body: nil,
		Code: -1,
	}

	uri, err := c.updateURI()

	if err != nil {
		return response, err
	}

	res := uri + id
	fmt.Printf("Fetching from URI <%s>\n", res)
	req, err := http.NewRequest(http.MethodGet, res, nil)

	if err != nil {
		log.Fatal("Error reading request. ", err)
		return response, err
	}

	req.Header.Add("Content-type", c.ContentType)
	client := &http.Client{Timeout: checkTimeout(timeout)}
	resp, err := client.Do(req)

	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}

	response.Body = body
	response.Code = resp.StatusCode

	return response, nil
}

// Default value checking methods

func checkTimeout(timeout time.Duration) time.Duration {
	if timeout == 0 {
		return defaultTimeout
	}
	return timeout
}

func (p *parameters) checkContentType() string {
	base := p.ContentType
	if p.ContentType == "" {
		base = defaultContentType
	}
	return base
}

func (c *Client) checkProtocol(protocol string) string {
	if protocol == "" {
		return defaultProtocol
	}
	return protocol
}

func (c *Client) checkPort(port string) string {
	if port == "" {
		return defaultPort
	}
	return port
}

// NewClient constructor with default values check
func NewClient(host string, port string, protocol string, p parameters) (*Client, error) {
	var client Client

	if host == "" {
		return &client, fmt.Errorf("%s %w", "Host", ErrParamNotSet)
	}

	err := client.updateParameters(p)

	if err != nil {
		return &client, err
	}

	client = Client{
		host:       host,
		port:       client.checkPort(port),
		protocol:   client.checkProtocol(protocol),
		parameters: p,
	}

	// Check for parameters default value
	client.ContentType = p.checkContentType()

	return &client, nil
}

// GetObjID is a support function extracting <id> out of JSON data
// GetObjID input parameter is byte array data from API response
// Method handels byte array data as unstructured data by leveraging generic interface map
// Method returns id as string value. In case of error empty string returned
func GetObjID(input []byte) (string, error) {
	var id string
	// Declared array map of string with empty interface (unstructured input data)
	// which will hold the value of the parsed json. Parsing embedded object <id> in JSON and return.
	var raw map[string]interface{}
	err := json.Unmarshal(input, &raw)

	if err != nil {
		log.Fatal(err)
		return id, err
	}

	raw = raw["data"].(map[string]interface{})
	id = raw["id"].(string)

	return id, nil
}

// GetObjVersion is a support function extracting <version> out of JSON data
// GetObjVersion input parameter is byte array data from API response
// Method handels byte array data as unstructured data by leveraging generic interface map
// Method returns version as integer value. In case of error version < 0
func GetObjVersion(input []byte) (int, error) {
	// Declared array map of string with empty interface (unstructured input data)
	// which will hold the value of the parsed json. Parse embedded object <version> in JSON and return.
	var raw map[string]interface{}
	err := json.Unmarshal(input, &raw)

	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	raw = raw["data"].(map[string]interface{})
	// Converting version type from float to integer.
	// Investigate why is value mapped as float when using empty interface?
	version := int(raw["version"].(float64))

	return version, nil
}

// JSONPrettyPrint prints JSON data in a more readable way on terminal
func JSONPrettyPrint(body []byte) (string, error) {
	dst := &bytes.Buffer{}

	if err := json.Indent(dst, body, "", "  "); err != nil {
		return "", err
	}
	return dst.String(), nil
}
