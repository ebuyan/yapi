package server

import (
	"net/http"
	"os"
	"yapi/socket"

	"github.com/gorilla/mux"
)

type Http struct {
	Host string
	Handler
}

func NewHttp(conversation *socket.Conversation) Http {
	return Http{
		Host:    os.Getenv("HTTP_HOST"),
		Handler: NewHandler(conversation),
	}
}

func (Http *Http) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", Http.Handler.SendCommand).Methods("POST")
	r.HandleFunc("/", Http.Handler.GetLastState).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(Http.Host, nil)
}
