package glagol

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
)

type Device struct {
	id          string
	platform    string
	certificate string
	token       string
	host        string

	mu    sync.RWMutex
	state DeviceState

	refreshTokenHandler func(ctx context.Context, deviceId, platform string) (string, error)
}

func NewDevice(deviceId, platform, certificate string) *Device {
	return &Device{id: deviceId, platform: platform, certificate: certificate}
}

func (d *Device) GetId() string {
	return d.id
}

func (d *Device) GetState() []byte {
	d.mu.RLock()
	js, _ := json.Marshal(d.state)
	d.mu.RUnlock()
	return js
}

func (d *Device) SetState(state []byte) {
	d.mu.Lock()
	_ = json.Unmarshal(state, &d.state)
	d.mu.Unlock()
}

func (d *Device) GetHost() string {
	return d.host
}

func (d *Device) SetHost(ipAddr, port string) {
	host := url.URL{Scheme: "wss", Host: ipAddr + ":" + port, Path: "/"}
	d.host = host.String()
}

func (d *Device) SetRefreshTokenHandler(handler func(ctx context.Context, deviceId, platform string) (string, error)) {
	d.refreshTokenHandler = handler
}

func (d *Device) RefreshToken(ctx context.Context) (err error) {
	d.token, err = d.refreshTokenHandler(ctx, d.id, d.platform)
	return err
}

func (d *Device) GetToken() string {
	return d.token
}

func (d *Device) GetCertificate() string {
	return d.certificate
}

type DeviceState struct {
	State struct {
		PlayerState struct {
			Duration float64 `json:"duration"`
			HasPause bool    `json:"hasPause"`
			HasPlay  bool    `json:"hasPlay"`
			Progress float64 `json:"progress"`
			Subtitle string  `json:"subtitle"`
			Title    string  `json:"title"`
			Extra    struct {
				CoverURI  string `json:"coverURI"`
				StateType string `json:"stateType"`
			} `json:"extra"`
		} `json:"playerState"`
		Playing bool    `json:"playing"`
		Volume  float64 `json:"volume"`
	} `json:"state"`
}
