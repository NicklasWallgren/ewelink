package ewelink

import (
	"fmt"
	"math/rand"
)

var (
	iphoneModels = [3]string{"6,1", "6,2", "7,1"}
	romVersions  = [2]string{"10.0", "10.0.2"}
	imeiFormat   = "DF7425A0-%d-%d-9F5E-3BC9179E48FB"
)

// MobileDevice is the interface implemented by types that can deliver information about a mobile device.
type MobileDevice interface {
	Model() string
	Imei() string
	Os() string
	RomVersion() string
}

// IOSDevice holds the contextual information for a IOS device.
type IOSDevice struct {
	model      string
	imei       string
	os         string
	romVersion string
}

func (i IOSDevice) String() string {
	return fmt.Sprintf("%#v", i)
}

// Model returns the IOS model.
func (i IOSDevice) Model() string {
	return i.model
}

// Imei returns the IOS IMEI.
func (i IOSDevice) Imei() string {
	return i.imei
}

// Os returns the Os model.
func (i IOSDevice) Os() string {
	return i.os
}

// RomVersion returns the Rom version.
func (i IOSDevice) RomVersion() string {
	return i.romVersion
}

func newIOSDevice() *IOSDevice {
	return &IOSDevice{model: getRandomIphoneModel(), imei: getRandomImei(), os: "iOS", romVersion: getRandomRomVersion()}
}

func getRandomIphoneModel() string {
	return "iPhone" + iphoneModels[rand.Intn(len(iphoneModels))] // #nosec:G404
}

func getRandomRomVersion() string {
	return romVersions[rand.Intn(len(romVersions))] // #nosec:G404
}

func getRandomImei() string {
	const highest = 9999
	const lowest = 1000

	return fmt.Sprintf(imeiFormat,
		getRandomNumber(lowest, highest), getRandomNumber(lowest, highest))
}
