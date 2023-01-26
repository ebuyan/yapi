package main

import (
	"log"
	"os"
	"yapi/internal/glagol"
	"yapi/internal/server"
	"yapi/internal/socket"
	"yapi/internal/yandex"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalln("no .env.local file")
	}

	oauthToken, err := yandex.NewOAuthClient().GetToken()
	if err != nil {
		log.Fatalln(err)
	}

	client := glagol.NewClient(os.Getenv("DEVICE_ID"), oauthToken)
	station, err := client.GetDevice()
	if err != nil {
		log.Fatalln(err)
	}

	conversation := socket.NewConversation(station)
	soc := socket.NewSocket(conversation)

	if err = soc.Run(); err != nil {
		log.Fatalln(err)
	}

	srv := server.NewHttp(soc)
	if err = srv.Start(); err != nil {
		log.Fatalln(err)
	}
}
