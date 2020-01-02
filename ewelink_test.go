package ewelink

import (
	"context"
	"fmt"
	"testing"
)

func TestAuthentication(t *testing.T) {
	eweLink := New(NewConfiguration("eu"))

	authenticator := NewEmailAuthenticator("EMAIL", "PASSWORD")
	session, err := eweLink.Authenticate(context.Background(), authenticator)

	fmt.Println(session)
	fmt.Println(err)
}

func TestAuthenticationWithEmail(t *testing.T) {
	eweLink := New(NewConfiguration("eu"))

	eweLink.AuthenticateWithEmail(context.Background(), "EMAIL", "PASSWORD")
}
