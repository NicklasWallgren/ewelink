package ewelink

import "fmt"

const baseUrl = "https://%s-api.coolkit.cc:8080/api"

type configuration struct {
	Region string
	Url    string
}

func NewConfiguration(region string) *configuration {
	return &configuration{Region: region, Url: fmt.Sprintf(baseUrl, region)}
}
