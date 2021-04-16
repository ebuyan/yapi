package glagol

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Device struct {
	Id       string       `json:"id"`
	Platform string       `json:"platform"`
	Glagol   DeviceGlagol `json:"glagol"`

	Config DeviceConfig `json:"-"`
	Token  string       `json:"-"`
	State  DeviceState  `json:"-"`

	locked bool `json:"-"`
}

func (d *Device) GetState() []byte {
	js, _ := json.Marshal(d.State)
	return js
}

func (d *Device) SetState(state []byte) {
	s := DeviceState{}
	json.Unmarshal(state, &s)
	d.State = s
}

func (d *Device) GetHost() string {
	host := url.URL{Scheme: "wss", Host: d.Config.IpAddr + ":" + d.Config.Port, Path: "/"}
	return host.String()
}

func (d *Device) GetOrigin() http.Header {
	return http.Header{"Origin": {"http://yandex.ru/"}}
}

func (d *Device) GetToken() string {
	return d.Token
}

func (d *Device) GetSertificate() string {
	return d.Glagol.Security.ServerCertificate
}

func (d *Device) Locked() bool {
	return d.locked
}

func (d *Device) Lock() {
	d.locked = true
}

func (d *Device) Unlock() {
	d.locked = false
}

type DeviceConfig struct {
	Port   string
	IpAddr string
}

type DeviceGlagol struct {
	Security DeviceGlagolSecurity `json:"security"`
}

type DeviceGlagolSecurity struct {
	ServerCertificate string `json:"server_certificate"`
	ServerPrivateKey  string `json:"server_private_key"`
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
