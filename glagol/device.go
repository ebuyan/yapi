package glagol

type DeviceList []Device

type Device struct {
	Id       string       `json:"id"`
	Platform string       `json:"platform"`
	Glagol   DeviceGlagol `json:"glagol"`

	Config DeviceConfig `json:"-"`
	Token  string       `json:"-"`
	State  DeviceState  `json:"-"`

	locked bool `json:"-"`
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

func (d Device) Locked() bool {
	return d.locked
}

func (d Device) Lock() {
	d.locked = true
}

func (d Device) Unlock() {
	d.locked = false
}
