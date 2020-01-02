package ewelink

import "encoding/json"

type Encoder interface {
	encode(payload payloadInterface) ([]byte, error)
}

type jsonEncoder struct{}

func newJsonEncoder() Encoder {
	return &jsonEncoder{}
}

func (e jsonEncoder) encode(payload payloadInterface) ([]byte, error) {
	return json.Marshal(payload)
}
