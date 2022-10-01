package ewelink

import (
	"context"
	"fmt"
	"time"
)

// Authenticator is the interface implemented by types that authenticate a user.
type Authenticator interface {
	Authenticate(context context.Context, client Client, session *Session) error
}

type emailAuthenticator struct {
	Email    string
	Password string
}

// Authenticate authenticates using email as primary identifier.
func (e emailAuthenticator) Authenticate(context context.Context, client Client, session *Session) error {
	response, err := client.call(context, newAuthenticationRequest(
		buildEmailAuthenticationPayload(e.Email, e.Password, session), session))
	if err != nil {
		return fmt.Errorf("unable to authenticate using email. %w", err)
	}

	session.updateTokenAndResponse((response).(*AuthenticationResponse))

	return nil
}

type phoneNumberAuthenticator struct {
	PhoneNumber string
	Password    string
}

// Authenticate authenticates using phone number as the primary identifier.
func (p phoneNumberAuthenticator) Authenticate(context context.Context, client Client, session *Session) error {
	panic("implement me")
}

// NewEmailAuthenticator returns a new instance of 'NewEmailAuthenticator.
func NewEmailAuthenticator(email string, password string) Authenticator {
	return &emailAuthenticator{Email: email, Password: password}
}

// NewPhoneNumberAuthenticator returns a new instance of 'NewPhoneNumberAuthenticator.
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
		AppID:      session.Application.AppID,
		AppSecret:  session.Application.AppSecret,
		Imei:       session.MobileDevice.Imei(),
		Os:         session.MobileDevice.Os(),
		Model:      session.MobileDevice.Model(),
		RomVersion: session.MobileDevice.RomVersion(),
		AppVersion: session.Application.AppVersion,
	}
}
