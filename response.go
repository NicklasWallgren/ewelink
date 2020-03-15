package ewelink

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Response interface {
	Decode(payload io.ReadCloser) (Response, error)
	Envelope() Envelope
}

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

type httpResponse struct {
	httpResponseEnvelope
}

func (h *httpResponse) Decode(payload io.ReadCloser) (Response, error) {
	return h, decode(payload, h)
}

func (h httpResponse) Envelope() Envelope {
	return h
}

type websocketResponseEnvelope struct {
	Error    int    `json:"error"`
	Sequence string `json:"sequence"`
	Message  string `json:"reason"`
	DeviceId string `json:"deviceId"`
}

func (w websocketResponseEnvelope) Code() int {
	return w.Error
}

func (w websocketResponseEnvelope) Cause() string {
	return w.Message
}

type websocketResponse struct {
	websocketResponseEnvelope
}

func (w *websocketResponse) Decode(payload io.ReadCloser) (Response, error) {
	return w, decode(payload, w)
}

func (w websocketResponse) Envelope() Envelope {
	return w
}

type AuthenticationResponse struct {
	httpResponse
	At     string `json:"at"`
	Rt     string `json:"rt"`
	User   User   `json:"user"`
	Region string `json:"region"`
}

func (a AuthenticationResponse) String() string {
	return fmt.Sprintf("%#v", a)
}

func (a *AuthenticationResponse) Decode(payload io.ReadCloser) (Response, error) {
	return a, decode(payload, a)
}

type User struct {
	Id              string      `json:"_id"`
	Email           string      `json:"email"`
	Password        string      `json:"password"`
	AppId           string      `json:"appid"`
	CreatedAt       string      `json:"createdat"`
	ApiKey          string      `json:"apikey"`
	Online          bool        `json:"online"`
	OnlineTime      string      `json:"onlinetime"`
	Ip              string      `json:"ip"`
	Location        string      `json:"location"`
	Language        string      `json:"lang"`
	OfflineTime     string      `json:"offlinetime"`
	BindInformation interface{} `json:"bindInfos"`
	AppInformation  interface{} `json:"appInfos"`
	UserStatus      string      `json:"userstatus"`
	UnknownField01  int         `json:"__v"`
}

type AppInformation struct{}

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
	Deviceid  string        `json:"deviceid"`
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
		Startup       string   `json:"startup"`
		FwVersion     string   `json:"fwVersion"`
		Switch        string   `json:"switch"`
		Switches      []string `json:"switches"`
		ControlType   int      `json:"controlType"`
		PartnerApikey string   `json:"partnerApikey"`
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
	DevConfig    struct {
	} `json:"devConfig"`
	Uiid int `json:"uiid"`
	Tags struct {
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

func decode(subject io.ReadCloser, target interface{}) error {
	decoder := json.NewDecoder(subject)
	//decoder.DisallowUnknownFields()

	return decoder.Decode(&target)
}
