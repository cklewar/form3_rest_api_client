package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

const host string = "accountapi"
const accountID string = "de37f789-1604-7c5b-a1e5-3673ea9cc2db"

/*func TestTableNewClient(t *testing.T) {

	var tables = []Parameters{
		{
			BaseURI:  "/v1/organisation/",
			Resource: "accounts",
		},
		{
			Timeout:     0,
			BaseURI:     "/v1/organisation/",
			ContentType: "",
			Resource:    "accounts",
		},
	}

	for _, parameters := range tables {
		_, ok := NewClient("192.168.2.50", "", "", parameters)

		if ok != nil {
			t.Errorf("Wrong result. Got %t but wanted %t", ok, false)
		}
	}
}
*/

// Initialize parameters struct. Omit certain fields to test default value settings
var defaultParams Parameters = Parameters{
	BaseURI:  "/v1/organisation/",
	Resource: "accounts",
}

func TestNewClient(t *testing.T) {
	// Construct client with host param not set
	_, ok := NewClient("", "", "", defaultParams)
	if ok == nil {
		t.Errorf("Wrong result. Got nil but wanted error")
	}
}

func TestCreate(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	path := filepath.Join("../../examples/json/org_acc_create.json")
	createInputData, err := ioutil.ReadFile(path)

	if err != nil {
		t.Errorf("Error loading input data %v", err)
	}

	createResp, err := c.Create(createInputData, defaultTimeout)

	if err != nil {
		t.Errorf("Error creating data %v", err)
	}

	if createResp.Code != http.StatusCreated {
		t.Errorf("Error creating data. Got response code %v expected %v", createResp.Code, http.StatusCreated)
	}
}

func TestFetch(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	fetchResp, err := c.Fetch(accountID, defaultTimeout)

	if err != nil {
		t.Errorf("Error fetching data %v", err)
	}

	if fetchResp.Code != http.StatusOK {
		t.Errorf("Error fetching data. Got response code %v expected %v", fetchResp.Code, http.StatusOK)
	}
}

func TestDelete(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	deleteResp, err := c.Delete(accountID, 0, defaultTimeout)

	if err != nil {
		t.Errorf("Error deleting data %v", err)
	}

	if deleteResp.Code != http.StatusNoContent {
		t.Errorf("Error deleting data. Got response code %v expected %v", deleteResp.Code, http.StatusNoContent)
	}
}

func TestCreateAlreadyCreated(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	path := filepath.Join("../../examples/json/org_acc_create.json")
	createInputData, err := ioutil.ReadFile(path)

	if err != nil {
		t.Errorf("Error loading input data %v", err)
	}

	createResp, err := c.Create(createInputData, defaultTimeout)

	if err != nil {
		t.Errorf("Error creating data %v", err)
	}

	if createResp.Code != http.StatusCreated {
		t.Errorf("Error creating data. Got response code %v expected %v", createResp.Code, http.StatusCreated)
	}

	createResp, err = c.Create(createInputData, defaultTimeout)

	if err != nil {
		t.Errorf("Error creating data %v", err)
	}

	if createResp.Code != http.StatusConflict {
		t.Errorf("Error creating data. Got response code %v expected %v", createResp.Code, http.StatusCreated)
	}

	deleteResp, err := c.Delete(accountID, 0, defaultTimeout)

	if err != nil {
		t.Errorf("Error deleting data %v", err)
	}

	if deleteResp.Code != http.StatusNoContent {
		t.Errorf("Error deleting data. Got response code %v expected %v", deleteResp.Code, http.StatusNoContent)
	}
}

func TestFetchNotFound(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	fetchResp, err := c.Fetch(accountID, defaultTimeout)

	if err != nil {
		t.Errorf("Error fetching data %v", err)
	}

	if fetchResp.Code != http.StatusNotFound {
		t.Errorf("Fetching data ok. Got response code %v expected %v", fetchResp.Code, http.StatusNotFound)
	}
}

func TestUpdateParameters(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	err := c.UpdateBaseURI(defaultParams.BaseURI)

	if err != nil {
		t.Errorf("Error updating BaseURI %v", err)
	}

	err = c.UpdateResource(defaultParams.Resource)

	if err != nil {
		t.Errorf("Error updating Resource %v", err)
	}
}

func TestUpdateBaseURI(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	// Set mandatory BaseURI parameter to emtpty string
	err := c.UpdateBaseURI("")
	// Check error value
	if err == nil {
		t.Errorf("Error updating BaseURI. Got %v expected %v", err, fmt.Errorf("%q: %w", "BaseURI", ErrParamNotSet))
	}
}

func TestUpdateResource(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	// Set mandatory Resource parameter to emtpty string
	err := c.UpdateResource("")
	// Check error value
	if err == nil {
		t.Errorf("Error updating Resource. Got %v expected %v", err, fmt.Errorf("%q: %w", "Ressource", ErrParamNotSet))
	}
}

func TestCheckContentType(t *testing.T) {
	// Check empty string leads to default value set
	if defaultParams.ContentType == "" {
		ct := defaultParams.checkContentType()
		if ct != defaultContentType {
			t.Errorf("Wrong result. Got %s but wanted %s", ct, defaultContentType)
		}
	}
}

func TestCheckTimeout(t *testing.T) {
	// Construct client with default values
	var timeout time.Duration = checkTimeout(0)

	if timeout != defaultTimeout {
		t.Errorf("Wrong result. Got %s but wanted %s", timeout, defaultTimeout)
	}
}

func TestCheckPort(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	var port string = c.checkPort("")
	if port != defaultPort {
		t.Errorf("Wrong result. Got %s but wanted %s", port, defaultPort)
	}
}

func TestCheckProtocol(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	var protocol string = c.checkProtocol("")
	if protocol != defaultProtocol {
		t.Errorf("Wrong result. Got %s but wanted %s", protocol, defaultProtocol)
	}
}

func TestGetObjID(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)
	path := filepath.Join("../../examples/json/org_acc_create.json")
	createInputData, err := ioutil.ReadFile(path)

	if err != nil {
		t.Errorf("Error loading input data %v", err)
	}

	createResp, err := c.Create(createInputData, defaultTimeout)

	if err != nil {
		t.Errorf("Error creating data %v", err)
	}

	if createResp.Code != http.StatusCreated {
		t.Errorf("Error creating data. Got response code %v expected %v", createResp.Code, http.StatusCreated)
	}

	id, err := GetObjID(createResp.Body)

	if err != nil {
		t.Errorf("Error getting object id %v", err)
	}

	if id == "" {
		t.Errorf("Error getting object id %v", id)
	}

	deleteResp, err := c.Delete(accountID, 0, defaultTimeout)

	if err != nil {
		t.Errorf("Error deleting data %v", err)
	}

	if deleteResp.Code != http.StatusNoContent {
		t.Errorf("Error deleting data. Got response code %v expected %v", deleteResp.Code, http.StatusNoContent)
	}

}

func TestGetObjVersion(t *testing.T) {
	c, _ := NewClient(host, "", "", defaultParams)
	path := filepath.Join("../../examples/json/org_acc_create.json")
	createInputData, err := ioutil.ReadFile(path)

	if err != nil {
		t.Errorf("Error loading input data %v", err)
	}

	createResp, err := c.Create(createInputData, defaultTimeout)

	if err != nil {
		t.Errorf("Error creating data %v", err)
	}

	if createResp.Code != http.StatusCreated {
		t.Errorf("Error creating data. Got response code %v expected %v", createResp.Code, http.StatusCreated)
	}

	version, err := GetObjVersion(createResp.Body)

	if err != nil {
		t.Errorf("Error getting object version %v", err)
	}

	if version == -1 {
		t.Errorf("Error getting object version %v", version)
	}

	deleteResp, err := c.Delete(accountID, 0, defaultTimeout)

	if err != nil {
		t.Errorf("Error deleting data %v", err)
	}

	if deleteResp.Code != http.StatusNoContent {
		t.Errorf("Error deleting data. Got response code %v expected %v", deleteResp.Code, http.StatusNoContent)
	}
}
