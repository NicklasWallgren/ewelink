package ewelink

import (
	"context"
	"fmt"
	"testing"
)

func TestAuthentication(t *testing.T) {
	eweLink := New()

	// authenticate using email
	session, err := eweLink.AuthenticateWithEmail(
		context.Background(), NewConfiguration("eu"), "EMAIL", "PASSWORD")

	if err != nil {
		panic(err)
	}

	// retrieve the list of registered devices
	devices, err := eweLink.GetDevices(context.Background(), session)

	// turn on the outlet(s) of the first device
	response, err := eweLink.SetDevicePowerState(context.Background(), session, &devices.Devicelist[0], true)

	fmt.Println(response)
	fmt.Println(err)
}

func TestAuthenticationWithEmail(t *testing.T) {
	eweLink := New()

	eweLink.AuthenticateWithEmail(context.Background(), NewConfiguration("eu"), "EMAIL", "PASSWORD")
}
