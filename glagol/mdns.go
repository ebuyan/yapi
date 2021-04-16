package glagol

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/mdns"
)

type mDNS struct {
	service string
}

func NewMDNS() mDNS {
	return mDNS{"_yandexio._tcp"}
}

func (m *mDNS) SetConfig(device *Device) (err error) {
	entriesCh := make(chan *mdns.ServiceEntry)
	defer close(entriesCh)

	go func() {
		for entry := range entriesCh {
			if device.Id == m.getDeviceId(entry) {
				device.Config.Done = true
				device.Config.IpAddr = entry.AddrV4.String()
				device.Config.Port = strconv.Itoa(entry.Port)
				log.Println("Found device on: " + device.GetHost())
				return
			}
		}
	}()
	err = mdns.Lookup(m.service, entriesCh)
	return
}

func (mdns *mDNS) getDeviceId(entry *mdns.ServiceEntry) (deviceId string) {
	for _, field := range entry.InfoFields {
		entryData := strings.Split(field, "=")
		if len(entryData) == 2 && entryData[0] == "deviceId" {
			deviceId = entryData[1]
			return
		}
	}
	return
}
