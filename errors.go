package ewelink

import "fmt"

type apiError struct {
	Code    int
	Message string
}

func (r apiError) Error() string {
	return fmt.Sprintf("%#v", r)
}

type wrongRegionError struct {
	apiError
} // 301

type authenticationError struct {
	apiError
} // 401

type requestError struct {
	apiError
} // 400, 406

type genericError struct {
	apiError
} // 500
