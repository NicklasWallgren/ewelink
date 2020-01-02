package ewelink

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// Client is the interface implemented by types that can invoke the ewelink API
type Client interface {
	// call is responsible for making the HTTP call against ewelink API
	call(request Request, context context.Context) (Response, error)
}

type client struct {
	client        *http.Client
	encoder       Encoder
	decoder       Decoder
	configuration *configuration
}

func newClient(configuration *configuration) *client {
	return &client{client: &http.Client{}, encoder: newJsonEncoder(), decoder: newJsonDecoder(), configuration: configuration}
}

func (c client) call(request Request, context context.Context) (Response, error) {
	encoded, err := c.encoder.encode(request.Payload())

	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(request.Method(), c.configuration.Url+"/"+request.Uri(), request.Query(), encoded, request.Headers(),
		request.IsToBeSigned())

	if err != nil {
		return nil, err
	}

	resp, err := c.request(req.WithContext(context))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return c.decoder.decode(request.Response(), resp)
}

// newRequest creates and prepares a instance of http request
func (c client) newRequest(method string, url string, query *url.Values, body []byte, headers *http.Header, isSigned bool) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(string(body)))

	addHeaders(req, headers)
	addQueryParameters(req, query)

	if isSigned {
		req.Header.Add("Authorization", "Sign "+calculateHash(body))
	}

	return req, err
}

func addHeaders(request *http.Request, headers *http.Header) {
	request.Header.Add("Content-Type", "application/json")

	if headers == nil {
		return
	}

	for key, values := range *headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
}

func addQueryParameters(request *http.Request, query *url.Values) {
	if query == nil {
		return
	}

	request.URL.RawQuery = query.Encode()
}

func (c client) request(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}
