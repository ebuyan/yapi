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
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalln("No .env.local file")
	}

	oauthClient := yandex.NewOAuthClient()
	oauthToken, err := oauthClient.GetToken()
	if err != nil {
		log.Fatalln(err)
	}

	client := glagol.NewClient(os.Getenv("DEVICE_ID"), oauthToken)
	station, err := client.GetDevice()
	if err != nil {
		log.Fatalln(err)
	}

	conversation := socket.NewConversation(station)
	socket := socket.NewSocket(conversation)
	err = socket.Run()
	if err != nil {
		log.Fatalln(err)
	}

	server := server.NewHttp(&socket)
	server.Start()
}
