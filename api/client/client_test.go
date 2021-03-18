package client

import (
	"testing"
)

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

func TestFetch(t *testing.T) {
	c, _ := NewClient("192.168.2.50", "", "", defaultParams)

	c.Fetch()
}

func TestContentTypeBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient("192.168.2.50", "", "", defaultParams)

	if c.ContentType != defaultContentType {
		t.Errorf("Wrong result. Got %s but wanted %s", c.ContentType, defaultContentType)
	}
}

func TestTimeoutBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient("192.168.2.50", "", "", defaultParams)

	if c.Timeout != defaultTimeout {
		t.Errorf("Wrong result. Got %s but wanted %s", c.Timeout, defaultTimeout)
	}
}

func TestPortBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient("192.168.2.50", "", "", defaultParams)

	var port string = c.portBase("")
	if port != defaultPort {
		t.Errorf("Wrong result. Got %s but wanted %s", port, defaultPort)
	}
}

func TestProtocolBase(t *testing.T) {
	// Construct client with default values
	c, _ := NewClient("192.168.2.50", "", "", defaultParams)

	var protocol string = c.protocolBase("")
	if protocol != defaultProtocol {
		t.Errorf("Wrong result. Got %s but wanted %s", protocol, defaultProtocol)
	}
}
