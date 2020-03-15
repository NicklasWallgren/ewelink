# eWeLink API library

Golang library for the EweLink API. Control Sonoff/eWeLink smart devices.

[![Build Status](https://travis-ci.org/NicklasWallgren/ewelink.svg?branch=master)](https://travis-ci.org/NicklasWallgren/ewelink)
[![Go Report Card](https://goreportcard.com/badge/github.com/stretchr/testify)](https://goreportcard.com/report/github.com/NicklasWallgren/ewelink)
[![GoDoc](https://godoc.org/github.com/NicklasWallgren/ewelink?status.svg)](https://godoc.org/github.com/NicklasWallgren/ewelink) 

Check out the API Documentation http://godoc.org/github.com/NicklasWallgren/ewelink

# Installation
The library can be installed through `go get` 
```bash
go get github.com/NicklasWallgren/ewelink
```

# Supported versions
We support the two major Go versions, which are 1.12 and 1.13 at the moment.

# Features
- Retrieve devices
- Turn on/off devices
- Get power consumption [TODO]
- Listen to device events

# Examples 

## Initiate payment request
```go
import (
    "context"
    "fmt"
    "github.com/NicklasWallgren/ewelink"
)

instance := ewelink.New()

// authenticate using email
session, err := instance.AuthenticateWithEmail(
	context.Background(), ewelink.NewConfiguration("REGION"), "EMAIL", "PASSWORD")

// retrieve the list of registered devices
devices, err := instance.GetDevices(context.Background(), session)

// turn on the outlet(s) of the first device
response, err := instance.SetDevicePowerState(context.Background(), session, &devices.Devicelist[0], true)

fmt.Println(response)
fmt.Println(err)
```

## Unit tests
```bash
go test -v -race $(go list ./... | grep -v vendor)
```

## Contributing
  - Fork it!
  - Create your feature branch: `git checkout -b my-new-feature`
  - Commit your changes: `git commit -am 'Useful information about your new features'`
  - Push to the branch: `git push origin my-new-feature`
  - Submit a pull request

## Contributors
  - [Nicklas Wallgren](https://github.com/NicklasWallgren)
  - [All Contributors][link-contributors]

[link-contributors]: ../../contributors