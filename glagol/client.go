package glagol

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	config     StationConfig
	oAuthToken string
	baseUrl    string
}

func NewAPIClient(OAuthToken string) APIClient {
	return APIClient{
		oAuthToken: OAuthToken,
		config:     NewStationConfig(),
		baseUrl:    "https://quasar.yandex.net/glagol",
	}
}

func (api APIClient) GetLocalStation() *Device {
	devices, err := api.getDeviceList()
	if err != nil {
		panic(err)
	}

	device, err := api.discoverDevices(devices)
	if err != nil {
		panic(err)
	}

	return device
}

func (api APIClient) getDeviceList() (DeviceList, error) {
	responseBody, err := api.sendRequest("device_list")
	if err != nil {
		return nil, err
	}

	response := DeviceListSuccessfulResponse{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Devices) == 0 {
		return nil, errors.New("No devices found at account")
	}

	return response.Devices, nil
}

func (api APIClient) discoverDevices(devices DeviceList) (*Device, error) {
	for _, device := range devices {
		if device.Id == api.config.StationId {
			token, err := api.getJwtTokenForDevice(device)
			if err != nil {
				return nil, err
			}
			device.Token = token
			device.Discovery = DeviceLocalDiscovery{
				Discovered:   true,
				LocalAddress: api.config.StationAddr,
				LocalPort:    api.config.StationPort,
			}
			return device, nil
		}
	}

	return nil, errors.New("No station found in local network")
}

func (api APIClient) getJwtTokenForDevice(device *Device) (string, error) {
	responseBody, err := api.sendRequest("token?device_id=" + device.Id + "&platform=" + device.Platform)
	if err != nil {
		return "", err
	}

	response := TokenSuccessfulResponse{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (api APIClient) sendRequest(endPoint string) ([]byte, error) {
	url := api.baseUrl + "/" + endPoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Oauth "+api.oAuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type DeviceListSuccessfulResponse struct {
	Devices DeviceList `json:"devices"`
	Status  string     `json:"status"`
}

type TokenSuccessfulResponse struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}
