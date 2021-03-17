package client

import (
	"testing"
)

// Initialize parameters struct. Omit certain fields to test default value settings
var defaultParams Parameters = Parameters{
	BaseURI:  "/v1/organisation/",
	Resource: "accounts",
}

func TestContentTypeBase(t *testing.T) {
	// Construct client with default values
	_, c := NewClient("192.168.2.50", "", "", defaultParams)

	if c.ContentType != defaultContentType {
		t.Errorf("Wrong result. Got %s but wanted %s", c.ContentType, defaultContentType)
	}
}

func TestTimeoutBase(t *testing.T) {
	// Construct client with default values
	_, c := NewClient("192.168.2.50", "", "", defaultParams)

	if c.Timeout != defaultTimeout {
		t.Errorf("Wrong result. Got %s but wanted %s", c.Timeout, defaultTimeout)
	}
}

func TestPortBase(t *testing.T) {
	// Construct client with default values
	_, c := NewClient("192.168.2.50", "", "", defaultParams)

	var port string = c.portBase("")
	if port != defaultPort {
		t.Errorf("Wrong result. Got %s but wanted %s", port, defaultPort)
	}
}

func TestProtocolBase(t *testing.T) {
	// Construct client with default values
	_, c := NewClient("192.168.2.50", "", "", defaultParams)

	var protocol string = c.protocolBase("")
	if protocol != defaultProtocol {
		t.Errorf("Wrong result. Got %s but wanted %s", protocol, defaultProtocol)
	}
}
