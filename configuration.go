package ewelink

import (
	"fmt"
	"net/url"
)

const (
	baseURL          = "https://%s-api.coolkit.cc:8080/api"
	websocketHost    = "%s-pconnect3.coolkit.cc:8080"
	defaultAppID     = "YzfeftUVcZ6twZw1OoVKPRFYTrGEg01Q"
	defaultAppSecret = "4G91qSoboqYO4Y0XJ0LPPKIsq8reHdfa" // nolint:gosec
)

// Configuration contains the configuration specific fields.
type Configuration struct {
	// The user account region
	Region       string
	APIURL       string
	WebsocketURL *url.URL
	AppID        string
	AppSecret    string
}

func (c Configuration) String() string {
	return fmt.Sprintf("%#v", c)
}

// ConfigurationOptionFunc option definition.
type ConfigurationOptionFunc func(c *Configuration)

// WithAppID option func to provide custom AppID.
func WithAppID(appID string) ConfigurationOptionFunc {
	return func(c *Configuration) {
		c.AppID = appID
	}
}

// WithAppSecret option func to provide custom AppSecret.
func WithAppSecret(appSecret string) ConfigurationOptionFunc {
	return func(c *Configuration) {
		c.AppSecret = appSecret
	}
}

// NewConfiguration creates a new Configuration.
func NewConfiguration(region string, optionFunc ...ConfigurationOptionFunc) *Configuration {
	configuration := &Configuration{
		Region: region, APIURL: fmt.Sprintf(baseURL, region),
		WebsocketURL: &url.URL{Scheme: websocketScheme, Host: fmt.Sprintf(websocketHost, region), Path: websocketURI},
		AppID:        defaultAppID,
		AppSecret:    defaultAppSecret,
	}

	for _, configurationOptionFunc := range optionFunc {
		configurationOptionFunc(configuration)
	}

	return configuration
}
