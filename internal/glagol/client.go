package glagol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultUrl       = "https://quasar.yandex.net/glagol"
	deviceListAction = "device_list"
	tokenAction      = "token"
)

type Client struct {
	deviceId string
	token    string
	baseUrl  string
	client   *http.Client
}

func NewClient(baseUrl, deviceId, token string) Client {
	if baseUrl == "" {
		baseUrl = defaultUrl
	}
	return Client{
		deviceId: deviceId,
		token:    token,
		baseUrl:  baseUrl,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (g *Client) GetDevice(ctx context.Context, ipAddr, port string) (device *Device, err error) {
	devices, err := g.getDeviceList(ctx)
	if err != nil {
		return
	}

	deviceResp, err := g.discoverDevices(devices)
	if err != nil {
		return
	}

	device = NewDevice(deviceResp.Id, deviceResp.Platform, deviceResp.Glagol.Security.ServerCertificate)
	device.SetHost(ipAddr, port)
	device.SetRefreshTokenHandler(g.getJwtTokenForDevice)
	if err = device.RefreshToken(ctx); err != nil {
		return nil, err
	}
	return device, nil
}

func (g *Client) getDeviceList(ctx context.Context) ([]DeviceResponse, error) {
	responseBody, err := g.sendRequest(ctx, deviceListAction)
	if err != nil {
		return nil, err
	}

	var response DeviceListResponse
	_ = json.Unmarshal(responseBody, &response)
	list := response.Devices
	if len(list) == 0 {
		return nil, errors.New("no devices found at account")
	}
	return list, nil
}

func (g *Client) discoverDevices(devices []DeviceResponse) (device DeviceResponse, err error) {
	for _, device = range devices {
		if device.Id == g.deviceId {
			return
		}
	}
	return device, errors.New("no station found in local network")
}

func (g *Client) getJwtTokenForDevice(ctx context.Context, deviceId, platform string) (token string, err error) {
	responseBody, err := g.sendRequest(ctx, fmt.Sprintf("%s?device_id=%s&platform=%s", tokenAction, deviceId, platform))
	if err != nil {
		return
	}
	response := TokenResponse{}
	_ = json.Unmarshal(responseBody, &response)
	return response.Token, nil
}

func (g *Client) sendRequest(ctx context.Context, endPoint string) (response []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseUrl+"/"+endPoint, http.NoBody)
	req.Header.Set("Authorization", "OAuth "+g.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode > http.StatusOK {
		return nil, errors.New(fmt.Sprintf("bad response %d", resp.StatusCode))
	}

	defer func() { _ = resp.Body.Close() }()
	response, err = ioutil.ReadAll(resp.Body)
	return
}

type DeviceListResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

type DeviceResponse struct {
	Id       string `json:"id"`
	Platform string `json:"platform"`
	Glagol   struct {
		Security struct {
			ServerCertificate string `json:"server_certificate"`
			ServerPrivateKey  string `json:"server_private_key"`
		} `json:"security"`
	} `json:"glagol"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
