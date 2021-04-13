package glagol

type DeviceResponse struct {
	State ResponseState `json:"state"`
}

type ResponseState struct {
	PlayerState ResponsePlayerState `json:"playerState"`
	Playing     bool                `json:"playing"`
	Volume      float64             `json:"volume"`
}

type ResponsePlayerState struct {
	Duration float64       `json:"duration"`
	Extra    ResponseExtra `json:"extra"`
	HasPause bool          `json:"hasPause"`
	HasPlay  bool          `json:"hasPlay"`
	Progress float64       `json:"progress"`
	Subtitle string        `json:"subtitle"`
	Title    string        `json:"title"`
}

type ResponseExtra struct {
	CoverURI string `json:"coverURI"`
}
