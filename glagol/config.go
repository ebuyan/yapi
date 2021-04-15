package glagol

import (
	"net/http"
	"net/url"
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

func (d DeviceConfig) GetHost() string {
	host := url.URL{Scheme: "wss", Host: d.IpAddr + ":" + d.Port, Path: "/"}
	return host.String()
}

func (d DeviceConfig) GetHeaderOrigin() http.Header {
	return http.Header{"Origin": {"http://yandex.ru/"}}
}
