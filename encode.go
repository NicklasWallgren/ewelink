package ewelink

import "encoding/json"

type encoder interface {
	encode(payload interface{}) ([]byte, error)
}

type jsonEncoder struct{}

func newJSONEncoder() encoder {
	return &jsonEncoder{}
}

func (e jsonEncoder) encode(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}
