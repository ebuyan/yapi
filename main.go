package main

import (
	"log"
	"yapi/glagol"
	"yapi/server"
	"yapi/socket"
	"yapi/yandex-oauth"

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

	client := glagol.NewGlagolClient(oauthToken)
	station, err := client.GetDevice()
	if err != nil {
		log.Fatalln(err)
	}

	socket := socket.NewSocket(&station)
	err = socket.Run()
	if err != nil {
		log.Fatalln(err)
	}

	server := server.NewHttp(&socket)
	server.Start()
}
