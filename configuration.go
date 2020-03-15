package ewelink

import "fmt"

const (
	baseUrl       = "https://%s-api.coolkit.cc:8080/api"
	websocketHost = "%s-pconnect3.coolkit.cc:8080"
)

type configuration struct {
	Region        string // user specified region
	ApiUrl        string // http api url
	WebsocketHost string // host or host:port
}

func NewConfiguration(region string) *configuration {
	return &configuration{Region: region, ApiUrl: fmt.Sprintf(baseUrl, region), WebsocketHost: fmt.Sprintf(websocketHost, region)}
}
