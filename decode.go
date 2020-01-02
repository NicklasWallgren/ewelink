package ewelink

import (
	"net/http"
)

type Decoder interface {
	decode(subject Response, response *http.Response) (Response, error)
}

type jsonDecoder struct{}

func newJsonDecoder() Decoder {
	return &jsonDecoder{}
}

func (j jsonDecoder) decode(subject Response, response *http.Response) (Response, error) {
	// Decode the response into the expected response type
	decoded, err := subject.Decode(response)

	if err != nil {
		return nil, err
	}

	// Check whether we encountered an API error
	if decoded.Envelope().Code > 0 {
		return nil, j.decodeAsApiError(decoded)
	}

	return decoded, nil
}

func (j jsonDecoder) decodeAsApiError(response Response) error {
	envelope := response.Envelope()

	return &apiError{Code: envelope.Code, Message: envelope.Message}
}
