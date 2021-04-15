package yandex

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type OAuthClient struct {
	baseUrl string
	body    OAuthRequestBody
	store   Store
}

func NewOAuthClient() OAuthClient {
	return OAuthClient{
		body:    NewOAuthRequestBody(),
		store:   NewStore(),
		baseUrl: "https://oauth.yandex.com/token",
	}
}

func (client OAuthClient) GetToken() (string, error) {
	ok, token := client.store.getToken()
	if !ok {
		resp, err := client.sendRequest()
		if err != nil {
			return "", err
		}
		if len(resp.CapchaKey) > 0 {
			return client.inputCaptcha(resp.CapchaUrl, resp.CapchaUrl)
		}
		token = resp.Token
		client.store.setToken(token)
	}
	return token, nil
}

func (client OAuthClient) sendRequest() (response OAuthTokenResponse, err error) {
	req, _ := http.NewRequest(http.MethodPost, client.baseUrl, bytes.NewBuffer(client.body.Get()))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response = OAuthTokenResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	if len(response.Error) > 0 {
		err = errors.New(response.Error)
		return
	}
	return
}

func (client OAuthClient) inputCaptcha(captchaUrl string, captchaKey string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please follow the link and enter the captcha value. " + captchaUrl)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if len(text) > 0 {
			client.body.captchaAnswer = text
			client.body.captchaKey = captchaKey
			return client.GetToken()
		}
	}
}

type OAuthTokenResponse struct {
	Token     string `json:"access_token"`
	CapchaKey string `json:"x_captcha_key"`
	CapchaUrl string `json:"x_captcha_url"`
	Error     string `json:"error_description"`
}
