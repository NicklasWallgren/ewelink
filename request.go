package ewelink

import (
	"net/http"
	"net/url"
)

const uriLogin = "user/login"
const uriGetDevices = "user/device"

type Request interface {
	Method() string
	Uri() string
	Query() *url.Values
	Payload() payloadInterface
	Headers() *http.Header
	Response() Response
	IsToBeSigned() bool
}

type request struct {
	method   string
	uri      string
	query    *url.Values
	payload  payloadInterface
	headers  *http.Header
	response Response
	isSigned bool
}

func (r request) Method() string {
	return r.method
}

func (r request) Uri() string {
	return r.uri
}

func (r request) Query() *url.Values {
	return r.query
}

func (r request) Payload() payloadInterface {
	return r.payload
}

func (r request) Headers() *http.Header {
	return r.headers
}

func (r request) Response() Response {
	return r.response
}

func (r request) IsToBeSigned() bool {
	return r.isSigned
}

func newAuthenticationRequest(payload *emailAuthenticationPayload) Request {
	return &request{method: "POST", uri: uriLogin, payload: payload, response: &AuthenticationResponse{}, isSigned: true}
}

func newGetDevicesRequest(query *url.Values, token string) Request {
	headers := &http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	return &request{method: "GET", uri: uriGetDevices, headers: headers, query: query, response: &DevicesResponse{}}
}
