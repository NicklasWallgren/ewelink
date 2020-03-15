package ewelink

import (
	"context"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
)

// ewelink contains the validator and configuration context
type ewelink struct {
	validator       *validator.Validate
	client          Client
	websocketClient WebsocketClient
}

// New returns a new instance of 'ewelink'
func New() *ewelink {
	return &ewelink{client: newClient(), websocketClient: newWebsocketClient()}
}

// Authenticate - Authenticates a new user session using an authenticator
func (e ewelink) Authenticate(context context.Context, configuration *configuration, authenticator Authenticator, options ...SessionOption) (*Session, error) {
	var session = &Session{IOSDevice: newIOSDevice(), Application: newApplication(), Configuration: configuration}

	// Apply options if there are any, can overwrite defaults
	for _, option := range options {
		option(session)
	}

	if err := authenticator.Authenticate(context, e.client, session); err != nil {
		return nil, err
	}

	return session, nil
}

// AuthenticateWithEmail - Authenticates a new user session using email as identifier
func (e ewelink) AuthenticateWithEmail(context context.Context, configuration *configuration, email string, password string, options ...SessionOption) (*Session, error) {
	return e.Authenticate(context, configuration, NewEmailAuthenticator(email, password), options...)
}

// GetDevices - Returns information about the devices
func (e ewelink) GetDevices(ctx context.Context, session *Session) (*DevicesResponse, error) {
	request := newGetDevicesRequest(createDeviceQuery(session), session.AuthenticationToken, session)

	response, err := e.call(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*DevicesResponse), nil
}

// GetDevice - Returns information about a device
func (e ewelink) GetDevice(deviceId string) {
	panic("implement")
}

// SetDevicePowerState - Toggles the outlet(s) of a device
func (e ewelink) SetDevicePowerState(context context.Context, session *Session, device *Device, stateOn bool) (Response, error) {
	numberOfOutlets, err := getDeviceOutletsCount(strconv.Itoa(device.Uiid))

	if err != nil {
		return nil, err
	}

	request := newWebsocketRequest(
		createUpdateActionPayload(createUpdatePowerStateOfDeviceParameters(numberOfOutlets, stateOn), session.AuthenticationResponse.User.ApiKey, device.Deviceid))

	return e.websocketCall(context, session, request)
}

// SetDeviceOutletPowerState - Toggles an outlet of a device
// The outlet indices start at 0.
func (e ewelink) SetDeviceOutletPowerState(context context.Context, session *Session, device *Device, stateOn bool, outletIndex int) (Response, error) {
	numberOfOutlets, err := getDeviceOutletsCount(strconv.Itoa(device.Uiid))

	if err != nil {
		return nil, err
	}

	if outletIndex < 0 || outletIndex > numberOfOutlets-1 {
		return nil, fmt.Errorf("invalid outlet index. Number of outlets for device: %d", numberOfOutlets)
	}

	var parameters interface{}

	if numberOfOutlets == 1 {
		parameters = createUpdatePowerStateOfDeviceParameters(numberOfOutlets, stateOn)
	} else {
		parameters = createUpdatePowerStateOfOutletParameters(device, outletIndex, numberOfOutlets, stateOn)
	}

	request := newWebsocketRequest(
		createUpdateActionPayload(parameters, session.AuthenticationResponse.User.ApiKey, device.Deviceid))

	return e.websocketCall(context, session, request)
}

func (e ewelink) websocketCall(context context.Context, session *Session, request WebsocketRequest) (Response, error) {
	authenticationRequest := newWebsocketRequest(createAuthenticateActionPayload(session))

	// TODO, validate the payload of the request?

	// Always authenticate as a first step/request
	responses, err := e.websocketClient.call(context, []WebsocketRequest{authenticationRequest, request}, session)

	if err != nil {
		return nil, err
	}

	for i, response := range responses {
		if response.Error != nil {
			return nil, response.Error
		}

		// Skip the response of the authentication request
		if i == 0 {
			continue
		}

		return response.Response, nil
	}

	return nil, fmt.Errorf("websocket request is nil, not allowed")
}

func (e ewelink) call(context context.Context, request HttpRequest) (Response, error) {
	if err := e.validate(request); err != nil {
		return nil, err
	}

	return e.client.call(request, context)
}

func (e ewelink) validate(request HttpRequest) error {
	if request.Payload() == nil {
		return nil
	}

	// Validate the integrity of the payload
	if err := e.validator.Struct(request.Payload()); err != nil {
		return err
	}

	return nil
}
