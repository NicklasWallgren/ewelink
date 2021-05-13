package ewelink

import (
	"fmt"
)

const (
	wrongRegion         = 301
	authenticationError = 401
	invalidRequest      = 400
	notAcceptable       = 406
	internalError       = 500
)

// APIErrorCause corresponds to a typical api error cause.
type APIErrorCause string

// APIErrorCauses holds information about all known api error causes.
var APIErrorCauses = struct {
	WrongRegion         APIErrorCause
	AuthenticationError APIErrorCause
	InvalidRequest      APIErrorCause
	InternalError       APIErrorCause
	UnknownError        APIErrorCause
}{
	WrongRegion:         "Wrong region",
	AuthenticationError: "Authentication required",
	InvalidRequest:      "Invalid request",
	InternalError:       "Internal server error",
}

type apiError struct {
	Code    int
	Message string
	Cause   APIErrorCause
}

type wrongRegionError struct {
	apiError
	Region string
}

func (r apiError) Error() string {
	return fmt.Sprintf("%#v", r)
}
