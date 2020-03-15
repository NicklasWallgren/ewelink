package ewelink

import (
	"net/http"
	"net/url"
)

const uriLogin = "user/login"
const uriGetDevices = "user/device"

type HttpRequest interface {
	Method() string
	Uri() string
	Query() *url.Values
	Payload() payloadInterface
	Headers() *http.Header
	Response() Response
	IsToBeSigned() bool
	Session() *Session
}

type httpRequest struct {
	method   string
	uri      string
	query    *url.Values
	payload  payloadInterface
	headers  *http.Header
	response Response
	isSigned bool
	session  *Session
}

func (r httpRequest) Method() string {
	return r.method
}

func (r httpRequest) Uri() string {
	return r.uri
}

func (r httpRequest) Query() *url.Values {
	return r.query
}

func (r httpRequest) Payload() payloadInterface {
	return r.payload
}

func (r httpRequest) Headers() *http.Header {
	return r.headers
}

func (r httpRequest) Response() Response {
	return r.response
}

func (r httpRequest) IsToBeSigned() bool {
	return r.isSigned
}

func (r httpRequest) Session() *Session {
	return r.session
}

type WebsocketRequest interface {
	Payload() payloadInterface
	Response() Response
	Session() *Session
}

type websocketRequest struct {
	payload  payloadInterface
	response Response
	session  *Session
}

func (w websocketRequest) Response() Response {
	return w.response
}

func (w websocketRequest) Payload() payloadInterface {
	return w.payload
}

func (w websocketRequest) Session() *Session {
	return w.session
}

func newAuthenticationRequest(payload *emailAuthenticationPayload, session *Session) HttpRequest {
	return &httpRequest{method: "POST", uri: uriLogin, payload: payload, response: &AuthenticationResponse{}, isSigned: true, session: session}
}

func newGetDevicesRequest(query *url.Values, token string, session *Session) HttpRequest {
	headers := &http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	return &httpRequest{method: "GET", uri: uriGetDevices, headers: headers, query: query, response: &DevicesResponse{}, session: session}
}

func newWebsocketRequest(payloadInterface payloadInterface) WebsocketRequest {
	return &websocketRequest{payload: payloadInterface, response: &websocketResponse{}}
}

func createUpdatePowerStateOfDeviceParameters(numberOfOutlets int, stateOn bool) interface{} {
	if numberOfOutlets == 1 {
		return &DevicePowerStateParameters{PowerOn: stateOn}
	}

	var statePerOutlet []bool

	for i := range statePerOutlet {
		statePerOutlet[i] = stateOn
	}

	return &DeviceOutletPowerStateAction{PowerOn: statePerOutlet}
}

func createUpdatePowerStateOfOutletParameters(device *Device, outletIndex int, numberOfOutlets int, stateOn bool) interface{} {
	var statePerOutlet []bool

	for i := range statePerOutlet {
		if i == outletIndex {
			statePerOutlet[i] = stateOn

			continue
		}

		statePerOutlet[i] = device.Params.Switches[i] == "on"
	}

	parameters := &DeviceOutletPowerStateAction{PowerOn: statePerOutlet}

	return parameters
}