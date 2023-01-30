package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"yapi/internal/glagol"
	"yapi/internal/server"
	"yapi/internal/socket"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalln("no .env.local file")
	}

	oauthToken := os.Getenv("OAUTH_TOKEN")
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
