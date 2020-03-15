package ewelink

import (
	"io"
)

type ResponseDecoder interface {
	decode(subject Response, response io.ReadCloser) (Response, error)
}

type responseJsonDecoder struct{}

func newResponseJsonDecoder() ResponseDecoder {
	return &responseJsonDecoder{}
}

func (j responseJsonDecoder) decode(subject Response, response io.ReadCloser) (Response, error) {
	// Decode the response into the expected response type
	decoded, err := subject.Decode(response)

	if err != nil {
		return nil, err
	}

	// Check whether we encountered an API error
	if decoded.Envelope().Code() > 0 {
		return nil, j.decodeAsApiError(decoded)
	}

	return decoded, nil
}

func (j responseJsonDecoder) decodeAsApiError(response Response) error {
	envelope := response.Envelope()

	// TODO, Better handling of the different error causes, websocket vs http
	// Include the actual response?

	return &apiError{Code: envelope.Code(), Message: envelope.Cause()}
}
