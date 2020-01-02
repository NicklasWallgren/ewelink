package ewelink

import (
	"net/url"
	"time"
)

// payload holds the request fields to be delivered to the API
type payload struct {
	payloadInterface
}

// payloadInterface is the interface implemented by types that holds the fields to be delivered to the API
type payloadInterface interface{}

// emailAuthenticationPayload holds the required and optional fields of the payment request
type emailAuthenticationPayload struct {
	*payload
	Email      string `json:"email"`
	Password   string `json:"password"`
	Version    string `json:"version"`
	Ts         int64  `json:"ts"`
	Nonce      string `json:"nonce"`
	Appid      string `json:"appid"`
	Imei       string `json:"imei"`
	Os         string `json:"os"`
	Model      string `json:"model"`
	RomVersion string `json:"romversion"`
	AppVersion string `json:"appversion"`
}

func createDeviceQuery(session *Session) *url.Values {
	query := &url.Values{}
	query.Add("lang", session.AuthenticationResponse.User.Language)
	query.Add("getTags", "1")
	query.Add("version", version)
	query.Add("ts", string(time.Now().Unix()))
	query.Add("appid", appId)

	return query
}

type devicesQuery struct {
	*payload
	Language string `json:"lang"`
	Tags     int    `json:"getTags"`
	Version  string `json:"version"`
	Ts       int64  `json:"ts"`
	AppId    string `json:"appid"`
}
