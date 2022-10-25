package ewelink

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// Ewelink contains the validator and Configuration context.
type Ewelink struct {
	validator       *validator.Validate
	client          Client
	websocketClient WebsocketClient
}

// New returns a new instance of 'Ewelink'.
func New(options ...Option) *Ewelink {
	ewelink := &Ewelink{client: newClient(), websocketClient: newWebsocketClient()}

	// Apply options if there are any, can overwrite defaults
	for _, option := range options {
		option(ewelink)
	}

	return ewelink
}

// Authenticate - Authenticates a new user session using an authenticator.
func (e Ewelink) Authenticate(context context.Context, configuration *Configuration, authenticator Authenticator, options ...SessionOptionFunc) (*Session, error) {
	session := &Session{MobileDevice: newIOSDevice(), Application: newApplication(), Configuration: configuration}

	// newSessionWith

	// Apply options if there are any, can overwrite defaults
	for _, option := range options {
		option(session)
	}

	if err := authenticator.Authenticate(context, e.client, session); err != nil {
		return nil, fmt.Errorf("failed to authenticate %w", err)
	}

	return session, nil
}

// AuthenticateWithEmail - Authenticates a new user session using email as identifier.
func (e Ewelink) AuthenticateWithEmail(context context.Context, configuration *Configuration, email string, password string, options ...SessionOptionFunc) (*Session, error) {
	return e.Authenticate(context, configuration, NewEmailAuthenticator(email, password), options...)
}

// AuthenticateWithPhoneNumber - Authenticates a new user session using phoneNumber as identifier.
func (e Ewelink) AuthenticateWithPhoneNumber(context context.Context, configuration *Configuration, phoneNumber string, password string, options ...SessionOptionFunc) (*Session, error) {
	return e.Authenticate(context, configuration, NewPhoneNumberAuthenticator(phoneNumber, password), options...)
}

// GetDevices - Returns information about the devices.
func (e Ewelink) GetDevices(ctx context.Context, session *Session) (*DevicesResponse, error) {
	request := newGetDevicesRequest(createDeviceQuery(session), session)

	response, err := e.call(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*DevicesResponse), nil
}

// GetDevice - Returns information about a device.
func (e Ewelink) GetDevice(ctx context.Context, session *Session, deviceID string) (*Device, error) {
	devices, err := e.GetDevices(ctx, session)
	if err != nil {
		return nil, err
	}

	return devices.getDeviceByID(deviceID)
}

// SetDevicePowerState - Toggles the outlet(s) of a device.
func (e Ewelink) SetDevicePowerState(context context.Context, session *Session, device *Device, stateOn bool) (Response, error) {
	numberOfOutlets, err := device.outletCount()
	if err != nil {
		return nil, err
	}

	request := newWebsocketRequest(
		createUpdateActionPayload(
			createUpdatePowerStateOfDeviceParameters(numberOfOutlets, stateOn), session.User.APIKey, device.DeviceID), &SetDevicePowerStateResponse{})

	response, err := e.websocketCall(context, session, request)
	if err != nil {
		return nil, err
	}

	return response.(*SetDevicePowerStateResponse), nil
}

// SetDeviceOutletPowerState - Toggles an outlet of a device
// The outlet indices start at 0.
func (e Ewelink) SetDeviceOutletPowerState(context context.Context, session *Session, device *Device, stateOn bool, outletIndex int) (Response, error) {
	if err := device.validOutletIndice(outletIndex); err != nil {
		return nil, err
	}

	numberOfOutlets, err := device.outletCount()
	if err != nil {
		return nil, err
	}

	var parameters parameters // TODO action parameters

	if numberOfOutlets == 1 {
		parameters = createUpdatePowerStateOfDeviceParameters(numberOfOutlets, stateOn)
	} else {
		parameters = createUpdatePowerStateOfOutletParameters(device, outletIndex, numberOfOutlets, stateOn)
	}

	request := newWebsocketRequest(
		createUpdateActionPayload(parameters, session.User.APIKey, device.DeviceID), &SetDeviceOutletPowerStateResponse{})

	response, err := e.websocketCall(context, session, request)
	if err != nil {
		return nil, err
	}

	return response.(*SetDeviceOutletPowerStateResponse), nil
}

func (e Ewelink) websocketCall(context context.Context, session *Session, request WebsocketRequest) (Response, error) {
	authenticationRequest := newWebsocketRequest(createAuthenticateActionPayload(session), &websocketResponse{})

	// Always authenticate as a first step/request
	responses, err := e.websocketClient.call(context, []WebsocketRequest{authenticationRequest, request}, session)
	if err != nil {
		return nil, fmt.Errorf("failed to send websocket request %w", err)
	}

	for i, response := range responses {
		if response.Error != nil {
			return nil, response.Error // TODO, wrap error, apiError?
		}

		// Skip the response of the authentication request
		if i == 0 {
			continue
		}

		return response.Response, nil
	}

	return nil, fmt.Errorf("websocket request is nil, not allowed")
}

func (e Ewelink) call(context context.Context, request HTTPRequest) (Response, error) {
	if err := e.validate(request); err != nil {
		return nil, err
	}

	return e.client.call(context, request)
}

func (e Ewelink) validate(request HTTPRequest) error {
	if request.Payload() == nil {
		return nil
	}

	// Validate the integrity of the payload
	if err := e.validator.Struct(request.Payload()); err != nil {
		return fmt.Errorf("valdation of http request failed %w", err)
	}

	return nil
}

// Option option definition.
type Option func(evelink *Ewelink)

// Function to create Option func to set custom http client.
// nolint:deadcode,unused
func withHTTPClient(client *http.Client) Option {
	return func(ewelink *Ewelink) {
		ewelink.client.withHTTPClient(client)
	}
}
