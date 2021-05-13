package ewelink

import (
	"encoding/json"
)

const (
	on  = "on"
	off = "off"
)

// payload is the interface implemented by types that holds the fields to be delivered to the API.
type payload interface{}

type parameters interface{}

// emailAuthenticationPayload holds the required and optional fields of the payment httpRequest.
type emailAuthenticationPayload struct {
	payload    `json:"-"` // nolint:unused
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Version    string     `json:"version"`
	Ts         int64      `json:"ts"`
	Nonce      string     `json:"nonce"`
	AppID      string     `json:"appid"`
	Imei       string     `json:"imei"`
	Os         string     `json:"os"`
	Model      string     `json:"model"`
	RomVersion string     `json:"romversion"`
	AppVersion string     `json:"appversion"`
}

// DevicePowerStateParameters for devices with only one outlet.
type DevicePowerStateParameters struct {
	parameters `json:"-"` // nolint:unused
	PowerOn    bool       `json:"switch"`
}

// MarshalJSON returns a JSON encoded 'DevicePowerStateParameters'.
func (d DevicePowerStateParameters) MarshalJSON() ([]byte, error) {
	power := off

	if d.PowerOn {
		power = on
	}

	return json.Marshal(&struct {
		Switch string `json:"switch"`
	}{Switch: power})
}

// DeviceOutletPowerStateAction device with multiple outlets.
type DeviceOutletPowerStateAction struct {
	parameters `json:"-"` // nolint:unused
	PowerOn    []bool     `json:"switches"`
}

// MarshalJSON returns a JSON encoded 'DeviceOutletPowerStateAction'.
func (d DeviceOutletPowerStateAction) MarshalJSON() ([]byte, error) {
	outlets := make([]string, len(d.PowerOn))

	for i, v := range d.PowerOn {
		state := on

		// nolint:gosimple
		if v == false {
			state = off
		}

		outlets[i] = state
	}

	return json.Marshal(&struct {
		Switch []string `json:"switches"`
	}{Switch: outlets})
}

// ActionPayload struct.
type ActionPayload struct {
	Action     string     `json:"action"`
	UserAgent  string     `json:"userAgent"`
	Parameters parameters `json:"params"`
	APIKey     string     `json:"apikey"`
	DeviceID   string     `json:"deviceid"`
	Sequence   int64      `json:"sequence"`
}

// AuthenticateActionPayload struct.
type AuthenticateActionPayload struct {
	Action     string `json:"action"`
	APIKey     string `json:"apikey"`
	AppID      string `json:"appid"`
	UserAgent  string `json:"userAgent"`
	Version    string `json:"version"`
	Nonce      string `json:"nonce"`
	ApkVersion string `json:"apkVersion"`
	Os         string `json:"os"`
	At         string `json:"at"`
	Ts         string `json:"ts"`
	Model      string `json:"model"`
	RomVersion string `json:"romVersion"`
	Sequence   int64  `json:"sequence"`
}

func createAuthenticateActionPayload(session *Session) *AuthenticateActionPayload {
	return &AuthenticateActionPayload{
		Action:     "userOnline",
		UserAgent:  "app",
		Version:    session.Application.Version,
		Nonce:      generateNonce(),
		ApkVersion: session.Application.ApkVersion,
		Os:         session.MobileDevice.Os(),
		At:         session.AuthenticationToken,
		APIKey:     session.User.APIKey,
		AppID:      session.User.AppID,
		Ts:         "1",
		Model:      session.MobileDevice.Model(),
		RomVersion: session.MobileDevice.RomVersion(),
		Sequence:   1,
	}
}

func createUpdateActionPayload(parameters parameters, apiKey string, deviceID string) *ActionPayload {
	return &ActionPayload{
		Action:     "update",
		UserAgent:  "app",
		Parameters: parameters,
		APIKey:     apiKey,
		DeviceID:   deviceID,
		Sequence:   2,
	}
}
