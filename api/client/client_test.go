package client

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
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
	_, ok := NewClient("", "", "", defaultParams)
	if ok == nil {
		t.Errorf("Wrong result. Got nil but wanted error")
	}
}

func TestCreate(t *testing.T) {
	c, _ := NewClient(host, "", "", defaultParams)
	path := filepath.Join("../../examples/json/org_acc_create.json")
	createInputData, err := ioutil.ReadFile(path)

	if err != nil {
		t.Errorf("Error loading input data %v", err)
	}

	createResp, err := c.Create(createInputData)

	if err != nil {
		t.Errorf("Error creating data %v", err)
	}

	if createResp.Code != http.StatusCreated {
		t.Errorf("Error creating data. Got response code %v expected %v", createResp.Code, http.StatusCreated)
	}
}

func TestFetch(t *testing.T) {
	c, _ := NewClient(host, "", "", defaultParams)
	fetchResp, err := c.Fetch(accountID)

	if err != nil {
		t.Errorf("Error fetching data %v", err)
	}

	if fetchResp.Code != http.StatusOK {
		t.Errorf("Error fetching data. Got response code %v expected %v", fetchResp.Code, http.StatusOK)
	}
}

func TestDelete(t *testing.T) {
	c, _ := NewClient(host, "", "", defaultParams)

	deleteResp, err := c.Delete(accountID, 0)

	if err != nil {
		t.Errorf("Error deleting data %v", err)
	}

	if deleteResp.Code != http.StatusNoContent {
		t.Errorf("Error deleting data. Got response code %v expected %v", deleteResp.Code, http.StatusNoContent)
	}
}

func TestContentTypeBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	if c.ContentType != defaultContentType {
		t.Errorf("Wrong result. Got %s but wanted %s", c.ContentType, defaultContentType)
	}
}

func TestTimeoutBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	if c.Timeout != defaultTimeout {
		t.Errorf("Wrong result. Got %s but wanted %s", c.Timeout, defaultTimeout)
	}
}

func TestPortBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	var port string = c.portBase("")
	if port != defaultPort {
		t.Errorf("Wrong result. Got %s but wanted %s", port, defaultPort)
	}
}

func TestProtocolBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient(host, "", "", defaultParams)

	var protocol string = c.protocolBase("")
	if protocol != defaultProtocol {
		t.Errorf("Wrong result. Got %s but wanted %s", protocol, defaultProtocol)
	}
}
