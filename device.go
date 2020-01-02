package ewelink

import (
	"fmt"
	"math/rand"
)

var (
	iphoneModels = [...]string{"6,1", "6,2", "7,1"}
	romVersions  = [...]string{"10.0", "10.0.2"}
)

type Device struct {
	Model      string
	Imei       string
	Os         string
	RomVersion string
}

func newDevice() *Device {
	return &Device{Model: getRandomIphoneModel(), Imei: getRandomImei(), Os: "iOS", RomVersion: getRandomRomVersion()}
}

func getRandomIphoneModel() string {
	return "iPhone" + iphoneModels[rand.Intn(len(iphoneModels))]
}

func getRandomRomVersion() string {
	return romVersions[rand.Intn(len(romVersions))]
}

func getRandomImei() string {
	const highest = 9999
	const lowest = 1000

	const imeiFormat = "DF7425A0-%d-%d-9F5E-3BC9179E48FB"
	return fmt.Sprintf(imeiFormat,
		getRandomNumber(lowest, highest), getRandomNumber(lowest, highest))
}
