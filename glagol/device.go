package glagol

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Device struct {
	id          string
	platform    string
	certificate string
	token       string
	host        string
	discovered  bool

	State DeviceState

	refreshTokenHandler func(deviceId, platform string) (string, error)
}

func NewDevice(deviceId, platform, certificate string) *Device {
	return &Device{id: deviceId, platform: platform, certificate: certificate}
}

func (d *Device) GetId() string {
	return d.id
}

func (d *Device) GetState() []byte {
	js, _ := json.Marshal(d.State)
	return js
}

func (d *Device) SetState(state []byte) {
	json.Unmarshal(state, &d.State)
}

func (d *Device) GetHost() string {
	return d.host
}

func (d *Device) SetHost(ipAddr, port string) {
	host := url.URL{Scheme: "wss", Host: ipAddr + ":" + port, Path: "/"}
	d.host = host.String()
	d.discovered = true
}

func (d *Device) GetOrigin() http.Header {
	return http.Header{"Origin": {"http://yandex.ru/"}}
}

func (d *Device) SetRefreshTokenHandler(handler func(deviceId, platform string) (string, error)) {
	d.refreshTokenHandler = handler
}

func (d *Device) RefreshToken() (err error) {
	token, err := d.refreshTokenHandler(d.id, d.platform)
	d.token = token
	return
}

func (d *Device) GetToken() string {
	return d.token
}

func (d *Device) GetCertificate() string {
	return d.certificate
}

func (d *Device) Discovered() bool {
	return d.discovered
}

type DeviceState struct {
	State State `json:"state"`
}

type State struct {
	PlayerState PlayerState `json:"playerState"`
	Playing     bool        `json:"playing"`
	Volume      float64     `json:"volume"`
}

type PlayerState struct {
	Duration float64 `json:"duration"`
	Extra    Extra   `json:"extra"`
	HasPause bool    `json:"hasPause"`
	HasPlay  bool    `json:"hasPlay"`
	Progress float64 `json:"progress"`
	Subtitle string  `json:"subtitle"`
	Title    string  `json:"title"`
}

type Extra struct {
	CoverURI string `json:"coverURI"`
}
