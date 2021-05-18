package yandex

import (
	"fmt"
	"os"
)

type OAuthRequestBody struct {
	clientId      string
	clientSecret  string
	grantType     string
	username      string
	password      string
	captchaAnswer string
	captchaKey    string
}

func NewOAuthRequestBody() OAuthRequestBody {
	return OAuthRequestBody{
		clientId:     "23cabbbdc6cd418abb4b39c32c41195d", //yandex-music
		clientSecret: "53bc75238f0c4d08a118e51fe9203300",
		grantType:    "password",
		username:     os.Getenv("LOGIN"),
		password:     os.Getenv("PASSWORD"),
	}
}

func (body OAuthRequestBody) String() string {
	str := fmt.Sprintf(
		"grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		body.grantType,
		body.clientId,
		body.clientSecret,
		body.username,
		body.password,
	)
	if len(body.captchaKey) > 0 {
		str += fmt.Sprintf(
			"&x_captcha_answer=%s&x_captcha_key=%s",
			body.captchaAnswer,
			body.captchaKey,
		)
	}
	return str
}
