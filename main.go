package main

import (
	"yapi/glagol"
	"yapi/server"
	"yapi/socket"
	"yapi/yandex-oauth"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		panic("No .env.local file")
	}

	oauthClient := yandex.NewOAuthClient()
	oauthToken := oauthClient.GetToken()

	api := glagol.NewAPIClient(oauthToken)
	station := api.GetLocalStation()

	conversation := socket.NewConversation(station)
	go conversation.Run()

	server := server.NewHttp(&conversation)
	server.Start()
}
