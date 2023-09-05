package ewelink

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Response interface.
type Response interface {
	Decode(payload io.ReadCloser) (Response, error)
	Envelope() Envelope
}

// Envelope interface.
type Envelope interface {
	Code() int
	Cause() string
}

type httpResponseEnvelope struct {
	Error   int    `json:"error"`
	Message string `json:"msg"`
	Region  string `json:"region"`
}

func (r httpResponseEnvelope) Code() int {
	return r.Error
}

func (r httpResponseEnvelope) Cause() string {
	return r.Message
}

// isErroneous returns true if an error encountered, false otherwise.
// nolint:unused
func (r httpResponseEnvelope) isErroneous() bool {
	return r.Error > 0
}

type httpResponse struct {
	httpResponseEnvelope
}

func (h *httpResponse) Decode(payload io.ReadCloser) (Response, error) {
	return h, decode(payload, h)
}

func (h httpResponse) Envelope() Envelope {
	return &h
}

type websocketResponseEnvelope struct {
	Error    int    `json:"error"`
	Sequence string `json:"sequence"`
	Message  string `json:"reason"`
	DeviceID string `json:"deviceId"`
}

func (w websocketResponseEnvelope) Code() int {
	return w.Error
}

func (w websocketResponseEnvelope) Cause() string {
	return w.Message
}

// isErroneous returns true if an error encountered, false otherwise.
// nolint:unused
func (w websocketResponseEnvelope) isErroneous() bool {
	return w.Error > 0
}

type websocketResponse struct {
	websocketResponseEnvelope
}

func (w *websocketResponse) Decode(payload io.ReadCloser) (Response, error) {
	return w, decode(payload, w)
}

func (w websocketResponse) Envelope() Envelope {
	return &w
}

// AuthenticationResponse struct.
type AuthenticationResponse struct {
	httpResponse
	At   string `json:"at"`
	Rt   string `json:"rt"`
	User User   `json:"user"`
}

func (a AuthenticationResponse) String() string {
	return fmt.Sprintf("%#v", a)
}

func (a *AuthenticationResponse) Decode(payload io.ReadCloser) (Response, error) {
	return a, decode(payload, a)
}

// User struct.
type User struct {
	ID              string      `json:"_id"`
	Email           string      `json:"email"`
	Password        string      `json:"password"`
	AppID           string      `json:"appid"`
	CreatedAt       string      `json:"createdat"`
	APIKey          string      `json:"apikey"`
	Online          bool        `json:"online"`
	OnlineTime      string      `json:"onlinetime"`
	IP              string      `json:"ip"`
	Location        string      `json:"location"`
	Language        string      `json:"lang"`
	OfflineTime     string      `json:"offlinetime"`
	BindInformation interface{} `json:"bindInfos"`
	AppInformation  interface{} `json:"appInfos"`
	UserStatus      string      `json:"userstatus"`
	UnknownField01  int         `json:"__v"`
}

func (u User) String() string {
	return fmt.Sprintf("%#v", u)
}

// SetDevicePowerStateResponse struct.
type SetDevicePowerStateResponse struct {
	websocketResponse
}

func (r SetDevicePowerStateResponse) String() string {
	return fmt.Sprintf("%#v", r)
}

func (r *SetDevicePowerStateResponse) Decode(payload io.ReadCloser) (Response, error) {
	return r, decode(payload, r)
}

// SetDeviceOutletPowerStateResponse struct.
type SetDeviceOutletPowerStateResponse struct {
	websocketResponse
}

func (r SetDeviceOutletPowerStateResponse) String() string {
	return fmt.Sprintf("%#v", r)
}

func (r *SetDeviceOutletPowerStateResponse) Decode(payload io.ReadCloser) (Response, error) {
	return r, decode(payload, r)
}

// AppInformation struct.
type AppInformation struct{}

// Device struct.
// nolint:maligned
type Device struct {
	Settings struct {
		OpsNotify   int `json:"opsNotify"`
		OpsHistory  int `json:"opsHistory"`
		AlarmNotify int `json:"alarmNotify"`
	} `json:"settings"`
	Group     string        `json:"group"`
	Online    bool          `json:"online"`
	Groups    []interface{} `json:"groups"`
	DevGroups []interface{} `json:"devGroups"`
	ID        string        `json:"_id"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	DeviceID  string        `json:"deviceid"`
	Apikey    string        `json:"apikey"`
	Extra     struct {
		Extra struct {
			Description  string `json:"description"`
			BrandID      string `json:"brandId"`
			Apmac        string `json:"apmac"`
			Mac          string `json:"mac"`
			UI           string `json:"ui"`
			ModelInfo    string `json:"modelInfo"`
			Model        string `json:"model"`
			Manufacturer string `json:"manufacturer"`
			Uiid         int    `json:"uiid"`
			StaMac       string `json:"staMac"`
		} `json:"extra"`
		ID string `json:"_id"`
	} `json:"extra"`
	CreatedAt  time.Time `json:"createdAt"`
	V          int       `json:"__v"`
	OnlineTime time.Time `json:"onlineTime"`
	IP         string    `json:"ip"`
	Location   string    `json:"location"`
	Params     struct {
		Rssi   int    `json:"rssi"`
		StaMac string `json:"staMac"`
		Timers []struct {
			Enabled          int    `json:"enabled"`
			Type             string `json:"type"`
			At               string `json:"at"`
			MID              string `json:"mId"`
			CoolkitTimerType string `json:"coolkit_timer_type"`
			Do               struct {
				Switch string `json:"switch"`
			} `json:"do"`
		} `json:"timers"`
		Startup   string `json:"startup"`
		FwVersion string `json:"fwVersion"`
		Switch    string `json:"switch"`
		Switches  []struct {
			Switch string `json:"switch"`
			Outlet int    `json:"outlet"`
		} `json:"switches"`
		ControlType   int    `json:"controlType"`
		PartnerApikey string `json:"partnerApikey"`
		BindInfos     struct {
			Gaction []string `json:"gaction"`
		} `json:"bindInfos"`
	} `json:"params"`
	OfflineTime  time.Time     `json:"offlineTime"`
	DeviceStatus string        `json:"deviceStatus"`
	SharedTo     []interface{} `json:"sharedTo"`
	Devicekey    string        `json:"devicekey"`
	DeviceURL    string        `json:"deviceUrl"`
	BrandName    string        `json:"brandName"`
	ShowBrand    bool          `json:"showBrand"`
	BrandLogoURL string        `json:"brandLogoUrl"`
	ProductModel string        `json:"productModel"`
	DevConfig    struct{}      `json:"devConfig"`
	Uiid         int           `json:"uiid"`
	Tags         struct {
		DisableTimers []struct {
			Do struct {
				Switch string `json:"switch"`
			} `json:"do"`
			CoolkitTimerType string `json:"coolkit_timer_type"`
			MID              string `json:"mId"`
			At               string `json:"at"`
			Type             string `json:"type"`
			Enabled          int    `json:"enabled"`
		} `json:"disable_timers"`
	} `json:"tags,omitempty"`
}

func (d Device) outletCount() (int, error) {
	return getDeviceOutletsCount(strconv.Itoa(d.Uiid))
}

func (d Device) validOutletIndice(outletIndex int) error {
	numberOfOutlets, err := d.outletCount()
	if err != nil {
		return err
	}

	if outletIndex < 0 || outletIndex > numberOfOutlets-1 {
		return fmt.Errorf("invalid outlet index. Number of outlets for device: %d", numberOfOutlets)
	}

	return nil
}

func (d *Device) String() string {
	return fmt.Sprintf("%#v", d)
}

// DevicesResponse struct.
type DevicesResponse struct {
	httpResponse
	Devicelist []Device `json:"devicelist"`
}

func (d *DevicesResponse) Decode(payload io.ReadCloser) (Response, error) {
	return d, decode(payload, d)
}

func (d DevicesResponse) String() string {
	return fmt.Sprintf("%#v", d)
}

func (d DevicesResponse) getDeviceByID(deviceID string) (*Device, error) {
	for _, device := range d.Devicelist {
		if device.DeviceID == deviceID {
			return &device, nil
		}
	}

	return nil, fmt.Errorf("invalid device id provided %s", deviceID)
}

func decode(subject io.ReadCloser, target interface{}) error {
	decoder := json.NewDecoder(subject)

	return decoder.Decode(&target)
}
