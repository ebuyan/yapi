package socket

import "net/http"

type Device interface {
	GetState() []byte
	SetState(state []byte)

	GetHost() string
	GetOrigin() http.Header

	RefreshToken() error
	GetToken() string
	GetCertificate() string
}
