package ewelink

import (
	"fmt"
	"net/url"
)

const (
	baseURL       = "https://%s-api.coolkit.cc:8080/api"
	websocketHost = "%s-pconnect3.coolkit.cc:8080"
)

// Configuration contains the configuration specific fields.
type Configuration struct {
	// The user account region
	Region       string
	APIURL       string
	WebsocketURL *url.URL
}

func (c Configuration) String() string {
	return fmt.Sprintf("%#v", c)
}

// NewConfiguration creates a new Configuration.
func NewConfiguration(region string) *Configuration {
	return &Configuration{
		Region: region, APIURL: fmt.Sprintf(baseURL, region),
		WebsocketURL: &url.URL{Scheme: websocketScheme, Host: fmt.Sprintf(websocketHost, region), Path: websocketURI},
	}
}
