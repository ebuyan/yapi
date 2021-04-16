package glagol

import (
	"os"
)

type DeviceConfig struct {
	Port   string
	IpAddr string
	Id     string
}

func NewDeviceConfig() DeviceConfig {
	return DeviceConfig{
		Port:   os.Getenv("STATION_PORT"),
		IpAddr: os.Getenv("STATION_ADDR"),
		Id:     os.Getenv("STATION_ID"),
	}
}
