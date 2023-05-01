package socket

import (
	"context"
	"net/http"
)

type Socket interface {
	Run(ctx context.Context) (err error)
	Read(w http.ResponseWriter, r *http.Request)
	Write(w http.ResponseWriter, r *http.Request)
}

type Device interface {
	GetState() []byte
	SetState(state []byte)
	GetHost() string
	RefreshToken(ctx context.Context) error
	GetToken() string
	GetCertificate() string
}
