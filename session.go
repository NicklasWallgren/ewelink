package ewelink

import "fmt"

// Session holds the metadata
type Session struct {
	IOSDevice              *IOSDevice // TODO, Device interface?
	Application            *Application
	AuthenticationToken    string
	AuthenticationResponse *AuthenticationResponse
	Configuration          *configuration
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

// Function to create SessionOption func to set custom ios device
func withIOSDevice(device *IOSDevice) SessionOption {
	return func(subject *Session) {
		subject.IOSDevice = device
	}
}

// Function to create SessionOption func to set custom application
func withApplication(application *Application) SessionOption {
	return func(subject *Session) {
		subject.Application = application
	}
}
