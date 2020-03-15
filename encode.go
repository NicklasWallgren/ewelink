package ewelink

import "encoding/json"

type Encoder interface {
	encode(payload interface{}) ([]byte, error)
}

type jsonEncoder struct{}

func newJsonEncoder() Encoder {
	return &jsonEncoder{}
}

func (e jsonEncoder) encode(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}
