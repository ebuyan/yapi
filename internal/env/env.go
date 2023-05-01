package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Env struct {
	OAuthToken string
	DeviceId   string
	HttpHost   string
	GlagolUrl  string
}

var Config Env

func init() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalln("no .env.local file")
	}
	Config.OAuthToken = os.Getenv("OAUTH_TOKEN")
	Config.DeviceId = os.Getenv("DEVICE_ID")
	Config.HttpHost = os.Getenv("HTTP_HOST")
	Config.GlagolUrl = os.Getenv("GLAGOL_URL")
}
