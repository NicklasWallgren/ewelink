package ewelink

import (
	"context"
	"time"
)

// Authenticator is the interface implemented by types that authenticate a user
type Authenticator interface {
	Authenticate(context context.Context, client Client, session *Session) error
}

type emailAuthenticator struct {
	Email    string
	Password string
}

//
func (e emailAuthenticator) Authenticate(context context.Context, client Client, session *Session) error {
	response, err := client.call(newAuthenticationRequest(
		buildEmailAuthenticationPayload(e.Email, e.Password, session)), context)

	if err != nil {
		return err
	}

	session.updateTokenAndResponse((response).(*AuthenticationResponse))

	return nil
}

type phoneNumberAuthenticator struct {
	PhoneNumber string
	Password    string
}

func (p phoneNumberAuthenticator) Authenticate(context context.Context, client Client, session *Session) error {
	panic("implement me")
}

// NewEmailAuthenticator returns a new instance of 'NewEmailAuthenticator'
func NewEmailAuthenticator(email string, password string) Authenticator {
	return &emailAuthenticator{Email: email, Password: password}
}

// NewPhoneNumberAuthenticator returns a new instance of 'NewPhoneNumberAuthenticator'
func NewPhoneNumberAuthenticator(phoneNumber string, password string) Authenticator {
	return &phoneNumberAuthenticator{PhoneNumber: phoneNumber, Password: password}
}

func buildEmailAuthenticationPayload(email string, password string, session *Session) *emailAuthenticationPayload {
	return &emailAuthenticationPayload{
		Email:      email,
		Password:   password,
		Version:    session.Application.Version,
		Ts:         time.Now().Unix(),
		Nonce:      generateNonce(),
		Appid:      session.Application.AppId,
		Imei:       session.Device.Imei,
		Os:         session.Device.Os,
		Model:      session.Device.Model,
		RomVersion: session.Device.RomVersion,
		AppVersion: session.Application.AppVersion,
	}
}
