package ewelink

import (
	"encoding/json"
	"net/url"
	"time"
)

// payload holds the httpRequest fields to be delivered to the API
type payload struct {
	payloadInterface
}

// payloadInterface is the interface implemented by types that holds the fields to be delivered to the API
type payloadInterface interface{}

// emailAuthenticationPayload holds the required and optional fields of the payment httpRequest
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

type DevicePowerStateParameters struct {
	PowerOn bool `json:"switch"`
}

func (d DevicePowerStateParameters) MarshalJSON() ([]byte, error) {
	power := "off"

	if d.PowerOn {
		power = "on"
	}

	return json.Marshal(&struct {
		Switch string `json:"switch"`
	}{Switch: power})
}

type DeviceOutletPowerStateAction struct {
	PowerOn []bool `json:"switches"`
}

func (d DeviceOutletPowerStateAction) MarshalJSON() ([]byte, error) {
	var outlets []string

	for i, v := range d.PowerOn {
		state := "on"

		if v == false {
			state = "off"
		}

		outlets[i] = state
	}

	return json.Marshal(&struct {
		Switch []string `json:"switches"`
	}{Switch: outlets})
}

type ActionPayload struct {
	Action     string      `json:"action"`
	UserAgent  string      `json:"userAgent"`
	Parameters interface{} `json:"params"`
	ApiKey     string      `json:"apikey"`
	DeviceId   string      `json:"deviceid"`
	Sequence   int64       `json:"sequence"`
}

type AuthenticateActionPayload struct {
	Action     string `json:"action"`
	UserAgent  string `json:"userAgent"`
	Version    string `json:"version"`
	Nonce      string `json:"nonce"`
	ApkVersion string `json:"apkVersion"`
	Os         string `json:"os"`
	At         string `json:"at"`
	ApiKey     string `json:"apikey"`
	Ts         string `json:"ts"`
	Model      string `json:"model"`
	RomVersion string `json:"romVersion"`
	Sequence   int64  `json:"sequence"`
}

func createAuthenticateActionPayload(session *Session) *AuthenticateActionPayload {

	// TODO, grab the information from the session metadata

	return &AuthenticateActionPayload{
		Action:     "userOnline",
		UserAgent:  "app",
		Version:    "6",
		Nonce:      generateNonce(),
		ApkVersion: "1.8",
		Os:         "ios",
		At:         session.AuthenticationToken,
		ApiKey:     session.AuthenticationResponse.User.ApiKey,
		Ts:         "1",
		Model:      "iPhone10,6",
		RomVersion: "11.1.2",
		Sequence:   1,
	}
}

func createUpdateActionPayload(parameters interface{}, apiKey string, deviceId string) *ActionPayload {
	return &ActionPayload{
		Action:     "update",
		UserAgent:  "app",
		Parameters: parameters,
		ApiKey:     apiKey,
		DeviceId:   deviceId,
		Sequence:   2,
	}
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
