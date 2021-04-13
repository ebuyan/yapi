package glagol

import (
	"os"
)

type StationConfig struct {
	StationPort string
	StationAddr string
	StationId   string
}

func NewStationConfig() StationConfig {
	return StationConfig{
		StationPort: os.Getenv("STATION_PORT"),
		StationAddr: os.Getenv("STATION_ADDR"),
		StationId:   os.Getenv("STATION_ID"),
	}
}
