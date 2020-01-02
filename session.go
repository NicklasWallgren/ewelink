package ewelink

import "fmt"

// Session holds the metadata
type Session struct {
	Device                 *Device
	Application            *Application
	AuthenticationToken    string
	AuthenticationResponse *AuthenticationResponse
}

func (s Session) String() string {
	return fmt.Sprintf("%#v", s)
}

func (s *Session) updateTokenAndResponse(response *AuthenticationResponse) {
	s.AuthenticationToken = response.At
	s.AuthenticationResponse = response
}

// Session option definition
type SessionOption func(session *Session)

// Function to create SessionOption func to set custom device
func withDevice(device *Device) SessionOption {
	return func(subject *Session) {
		subject.Device = device
	}
}

// Function to create SessionOption func to set custom application
func withApplication(application *Application) SessionOption {
	return func(subject *Session) {
		subject.Application = application
	}
}
