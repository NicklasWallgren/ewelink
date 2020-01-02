package ewelink

import (
	"context"
	"gopkg.in/go-playground/validator.v9"
)

// ewelink contains the validator and configuration context
type ewelink struct {
	validator *validator.Validate
	client    Client
}

// New returns a new instance of 'ewelink'
func New(configuration *configuration) *ewelink {
	return &ewelink{client: newClient(configuration)}
}

// Authenticate - Authenticates a new user session using an authenticator
func (e ewelink) Authenticate(context context.Context, authenticator Authenticator, options ...SessionOption) (*Session, error) {
	var session = &Session{Device: newDevice(), Application: newApplication()}

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
func (e ewelink) AuthenticateWithEmail(context context.Context, email string, password string, options ...SessionOption) (*Session, error) {
	return e.Authenticate(context, NewEmailAuthenticator(email, password), options...)
}

// GetDevices - Returns information about the devices
func (e ewelink) GetDevices(ctx context.Context, session *Session) (*DevicesResponse, error) {
	request := newGetDevicesRequest(createDeviceQuery(session), session.AuthenticationToken)

	response, err := e.call(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*DevicesResponse), nil
}

// GetDevice - Returns information about a device
func (e ewelink) GetDevice(id string) {

}

func (e ewelink) call(context context.Context, request Request) (Response, error) {
	if err := e.validate(request); err != nil {
		return nil, err
	}

	response, err := e.client.call(request, context)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (e ewelink) validate(request Request) error {
	if request.Payload() == nil {
		return nil
	}

	// Validate the integrity of the payload
	if err := e.validator.Struct(request.Payload()); err != nil {
		return err
	}

	return nil
}
