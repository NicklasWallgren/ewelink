package ewelink

import "math/rand"

const (
	appId   = "oeVkj2lYFGnJu5XUtWisfW4utiN4u9Mq"
	version = "6"
)

var applicationVersions = [...]string{"3.5.3", "3.5.4"}

//
type Application struct {
	AppVersion string
	Version    string
	AppId      string
}

func newApplication() *Application {
	return &Application{AppVersion: getRandomApplicationVersion(), Version: version, AppId: appId}
}

func getRandomApplicationVersion() string {
	return applicationVersions[rand.Intn(len(applicationVersions))]
}
