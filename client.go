package ewelink

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client is the interface implemented by types that can invoke the Ewelink API.
type Client interface {
	// call is responsible for making the HTTP call against Ewelink API
	call(context context.Context, request HTTPRequest) (Response, error)
	withHTTPClient(client *http.Client)
}

type client struct {
	client  *http.Client
	encoder encoder
	decoder responseDecoder
}

func newClient() *client {
	return &client{client: &http.Client{}, encoder: newJSONEncoder(), decoder: newJSONResponseDecoder()}
}

func (c client) call(context context.Context, request HTTPRequest) (Response, error) {
	encoded, err := c.encoder.encode(request.Payload())
	if err != nil {
		return nil, fmt.Errorf("could not encode http request payload. %w", err)
	}

	req, err := c.newRequest(context, request.Method(), buildRequestURL(request), request.Query(), encoded, request.Session().Configuration.AppSecret, request.Headers(), request.IsToBeSigned())
	if err != nil {
		return nil, fmt.Errorf("unable to process request. %w", err)
	}

	resp, err := c.request(req.WithContext(context))
	if err != nil {
		return nil, fmt.Errorf("error occured while calling api. %w", err)
	}

	defer resp.Body.Close() // nolint:errcheck

	return c.decoder.decode(request.Response(), resp.Body, resp.StatusCode)
}

func (c *client) withHTTPClient(client *http.Client) {
	c.client = client
}

// newRequest creates and prepares a instance of http http.httpRequest.
func (c client) newRequest(context context.Context, method string, url string, query *url.Values, body []byte, appSecret string, headers *http.Header, isSigned bool) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context, method, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("unable to create request %w", err)
	}

	addHeaders(req, headers)
	addQueryParameters(req, query)

	if isSigned {
		hashedBody, err := calculateHash(body, appSecret)
		if err != nil {
			return nil, fmt.Errorf("unable to calculated the hash of the request body %w", err)
		}

		req.Header.Add("Authorization", "Sign "+hashedBody)
	}

	return req, nil
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

func buildRequestURL(request HTTPRequest) string {
	return request.Session().Configuration.APIURL + "/" + request.URI()
}

func (c client) request(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}
