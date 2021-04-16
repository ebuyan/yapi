package glagol

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type GlagolClient struct {
	mdns     mDNS
	deviceId string
	token    string
	baseUrl  string
}

func NewGlagolClient(deviceId, token string) GlagolClient {
	return GlagolClient{
		mdns:     NewMDNS(),
		deviceId: deviceId,
		token:    token,
		baseUrl:  "https://quasar.yandex.net/glagol",
	}
}

func (g *GlagolClient) GetDevice() (device *Device, err error) {
	devices, err := g.getDeviceList()
	if err != nil {
		return
	}
	device, err = g.discoverDevices(devices)
	if err != nil {
		return
	}
	token, err := g.getJwtTokenForDevice(device)
	if err != nil {
		return
	}
	device.Token = token
	err = g.mdns.SetConfig(device)
	if err != nil {
		return
	}
	if !device.Config.Done {
		err = errors.New("Failed to resolve device IpAddr")
	}
	return
}

func (g *GlagolClient) getDeviceList() ([]*Device, error) {
	responseBody, err := g.sendRequest("device_list")
	if err != nil {
		return nil, err
	}
	response := DeviceListResponse{}
	json.Unmarshal(responseBody, &response)
	list := response.Devices
	if len(list) == 0 {
		err = errors.New("No devices found at account")
	}
	return list, err
}

func (g *GlagolClient) discoverDevices(devices []*Device) (device *Device, err error) {
	for _, device = range devices {
		if device.Id == g.deviceId {
			return
		}
	}
	err = errors.New("No station found in local network")
	return
}

func (api *GlagolClient) getJwtTokenForDevice(device *Device) (token string, err error) {
	responseBody, err := api.sendRequest("token?device_id=" + device.Id + "&platform=" + device.Platform)
	if err != nil {
		return
	}
	response := TokenResponse{}
	json.Unmarshal(responseBody, &response)
	token = response.Token
	return
}

func (g *GlagolClient) sendRequest(endPoint string) (response []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, g.baseUrl+"/"+endPoint, nil)
	req.Header.Set("Authorization", "Oauth "+g.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

type DeviceListResponse struct {
	Devices []*Device `json:"devices"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
