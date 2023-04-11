package main

import (
	"log"
	"yapi/internal/env"
	"yapi/internal/glagol"
	"yapi/internal/server"
	"yapi/internal/socket"
)

func main() {
	client := glagol.NewClient(env.Config.DeviceId, env.Config.OAuthToken)
	station, err := client.GetDevice()
	if err != nil {
		log.Fatalln(err)
	}

	conversation := socket.NewConversation(station)
	soc := socket.NewSocket(conversation)

	if err = soc.Run(); err != nil {
		log.Fatalln(err)
	}

	srv := server.NewHttp(soc, env.Config.HttpHost)
	if err = srv.Start(); err != nil {
		log.Fatalln(err)
	}
}
