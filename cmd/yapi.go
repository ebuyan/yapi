package main

import (
	"context"
	"log"
	"yapi/internal/env"
	"yapi/internal/glagol"
	"yapi/internal/server"
	"yapi/internal/socket"
	"yapi/pkg/mdns"
)

func main() {
	entry, err := mdns.Discover(env.Config.DeviceId, mdns.YandexServicePrefix)
	if err != nil {
		log.Fatalln(err)
	}

	client := glagol.NewClient(env.Config.GlagolUrl, env.Config.DeviceId, env.Config.OAuthToken)
	station, err := client.GetDevice(context.Background(), entry.IpAddr, entry.Port)
	if err != nil {
		log.Fatalln(err)
	}

	conversation := socket.NewConversation(station)
	soc := socket.NewSocket(conversation)

	if err = soc.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}

	srv := server.NewHttp(soc, env.Config.HttpHost)
	if err = srv.Start(); err != nil {
		log.Fatalln(err)
	}
}
