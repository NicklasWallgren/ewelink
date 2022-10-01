package ewelink

import (
	"fmt"
	"math/rand"
)

const (
	appID      = "YzfeftUVcZ6twZw1OoVKPRFYTrGEg01Q"
	appSecret  = "4G91qSoboqYO4Y0XJ0LPPKIsq8reHdfa"
	version    = "8"
	apkVersion = "1.8"
)

var applicationVersions = [2]string{"3.5.3", "3.5.4"}

// Application contains the application specific fields.
type Application struct {
	AppVersion string
	Version    string
	AppID      string
	AppSecret  string
	ApkVersion string
}

func (a Application) String() string {
	return fmt.Sprintf("%#v", a)
}

func NewApplication() *Application {
	return &Application{AppVersion: getRandomApplicationVersion(), Version: version, AppID: appID, AppSecret: appSecret, ApkVersion: apkVersion}
}

func getRandomApplicationVersion() string {
	return applicationVersions[rand.Intn(len(applicationVersions))] // #nosec:G404
}
