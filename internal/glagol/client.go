package glagol

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"yapi/pkg/mdns"
)

type GlagolClient struct {
	deviceId string
	token    string
	baseUrl  string
}

func NewGlagolClient(deviceId, token string) GlagolClient {
	return GlagolClient{
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
	deviceResp, err := g.discoverDevices(devices)
	if err != nil {
		return
	}
	device = NewDevice(deviceResp.Id, deviceResp.Platform, deviceResp.Glagol.Security.ServerCertificate)
	entry, err := mdns.Discover(device.id, "_yandexio._tcp")
	if err != nil {
		return
	}
	device.SetHost(entry.IpAddr, entry.Port)
	device.SetRefreshTokenHandler(g.getJwtTokenForDevice)
	err = device.RefreshToken()
	return
}

func (g *GlagolClient) getDeviceList() ([]DeviceResponse, error) {
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

func (g *GlagolClient) discoverDevices(devices []DeviceResponse) (device DeviceResponse, err error) {
	for _, device = range devices {
		if device.Id == g.deviceId {
			return
		}
	}
	err = errors.New("No station found in local network")
	return
}

func (g *GlagolClient) getJwtTokenForDevice(deviceId, platform string) (token string, err error) {
	responseBody, err := g.sendRequest("token?device_id=" + deviceId + "&platform=" + platform)
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
	Devices []DeviceResponse `json:"devices"`
}

type DeviceResponse struct {
	Id       string       `json:"id"`
	Platform string       `json:"platform"`
	Glagol   DeviceGlagol `json:"glagol"`
}

type DeviceGlagol struct {
	Security DeviceGlagolSecurity `json:"security"`
}

type DeviceGlagolSecurity struct {
	ServerCertificate string `json:"server_certificate"`
	ServerPrivateKey  string `json:"server_private_key"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
