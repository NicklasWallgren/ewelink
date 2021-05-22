package ewelink

import (
	"net/http"
	"net/url"
)

// HTTPRequest interface.
type HTTPRequest interface {
	Method() string
	URI() string
	Query() *url.Values
	Payload() payload
	Headers() *http.Header
	Response() Response
	IsToBeSigned() bool
	Session() *Session
}

type httpRequest struct {
	method   string
	uri      string
	query    *url.Values
	payload  payload
	headers  *http.Header
	response Response
	isSigned bool
	session  *Session
}

func (r httpRequest) Method() string {
	return r.method
}

func (r httpRequest) URI() string {
	return r.uri
}

func (r httpRequest) Query() *url.Values {
	return r.query
}

func (r httpRequest) Payload() payload {
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

// WebsocketRequest interface.
type WebsocketRequest interface {
	Payload() payload
	Response() Response
	Session() *Session
}

type websocketRequest struct {
	payload  payload
	response Response
	session  *Session
}

func (w websocketRequest) Response() Response {
	return w.response
}

func (w websocketRequest) Payload() payload {
	return w.payload
}

func (w websocketRequest) Session() *Session {
	return w.session
}

func newAuthenticationRequest(payload *emailAuthenticationPayload, session *Session) HTTPRequest {
	return &httpRequest{method: "POST", uri: "user/login", payload: payload, response: &AuthenticationResponse{}, isSigned: true, session: session}
}

func newGetDevicesRequest(query *url.Values, session *Session) HTTPRequest {
	headers := &http.Header{}
	headers.Add("Authorization", "Bearer "+session.AuthenticationToken)

	return &httpRequest{method: "GET", uri: "user/device", headers: headers, query: query, response: &DevicesResponse{}, session: session}
}

func newWebsocketRequest(payload payload, response Response) WebsocketRequest {
	return &websocketRequest{payload: payload, response: response}
}

func createUpdatePowerStateOfDeviceParameters(numberOfOutlets int, stateOn bool) parameters {
	if numberOfOutlets == 1 {
		return &DevicePowerStateParameters{PowerOn: stateOn}
	}

	statePerOutlet := make([]bool, numberOfOutlets)

	for i := range statePerOutlet {
		statePerOutlet[i] = stateOn
	}

	return &DeviceOutletPowerStateAction{PowerOn: statePerOutlet}
}

func createUpdatePowerStateOfOutletParameters(device *Device, outletIndex int, numberOfOutlets int, stateOn bool) parameters {
	statePerOutlet := make([]bool, numberOfOutlets)

	for i := range statePerOutlet {
		if i == outletIndex {
			statePerOutlet[i] = stateOn

			continue
		}

		statePerOutlet[i] = device.Params.Switches[i] == on
	}

	return &DeviceOutletPowerStateAction{PowerOn: statePerOutlet}
}
