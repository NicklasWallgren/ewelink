package main

import (
	"context"
	"fmt"

	"github.com/jacostaperu/ewelink"
)

func main() {
	instance := ewelink.New()

	authenticator := ewelink.NewEmailAuthenticator("EMAIL", "PASSWORD")

	session, err := instance.Authenticate(context.Background(), ewelink.NewConfiguration("REGION"), authenticator)

	fmt.Println(session)
	fmt.Println(err)
}
