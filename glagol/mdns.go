package glagol

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/mdns"
)

type MDNS struct {
	service string
}

func NewMDNS() MDNS {
	return MDNS{"_yandexio._tcp"}
}

func (m MDNS) SetIpAddrPort(device *Device) (err error) {
	entriesCh := make(chan *mdns.ServiceEntry)
	go func() {
		for entry := range entriesCh {
			if device.Id == m.GetDeviceId(entry) {
				device.Config.IpAddr = entry.AddrV4.String()
				device.Config.Port = strconv.Itoa(entry.Port)
				log.Println("Found device on: " + device.GetHost())
				return
			}
		}
		err = errors.New("Device not Found")
	}()
	mdns.Lookup(m.service, entriesCh)
	close(entriesCh)
	return
}

func (mdns MDNS) GetDeviceId(entry *mdns.ServiceEntry) (id string) {
	for _, field := range entry.InfoFields {
		entryData := strings.Split(field, "=")
		if len(entryData) == 2 && entryData[0] == "deviceId" {
			id = entryData[1]
			return
		}
	}
	return
}
