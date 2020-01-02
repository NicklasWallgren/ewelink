package main

import (
	"context"
	"ewelink"
	"fmt"
)

func main() {
	instance := ewelink.New(ewelink.NewConfiguration("eu"))
	authenticator := ewelink.NewEmailAuthenticator("EMAIL", "PASSWORD")
	session, err := instance.Authenticate(context.Background(), authenticator)

	fmt.Println(session)
	fmt.Println(err)
}


