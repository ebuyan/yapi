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
	"yapi/pkg/store"
)

type OAuthClient struct {
	baseUrl string
	body    OAuthRequestBody
	store   store.Store
}

func NewOAuthClient() OAuthClient {
	return OAuthClient{
		body:    NewOAuthRequestBody(),
		store:   store.NewStore(),
		baseUrl: "https://oauth.yandex.com/token",
	}
}

func (c OAuthClient) GetToken() (string, error) {
	ok, token := c.store.GetToken()
	if !ok {
		resp, err := c.sendRequest()
		if err != nil {
			return "", err
		}
		if len(resp.CaptchaKey) > 0 {
			return c.inputCaptcha(resp.CaptchaUrl, resp.CaptchaUrl)
		}
		token = resp.Token
		c.store.SetToken(token)
	}
	return token, nil
}

func (c OAuthClient) sendRequest() (response OAuthTokenResponse, err error) {
	req, _ := http.NewRequest(http.MethodPost, c.baseUrl, bytes.NewBuffer([]byte(c.body.String())))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response = OAuthTokenResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	if len(response.Error) > 0 {
		err = errors.New(response.Error)
	}
	return
}

func (c OAuthClient) inputCaptcha(captchaUrl string, captchaKey string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please follow the link and enter the captcha value. " + captchaUrl)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if len(text) > 0 {
			c.body.captchaAnswer = text
			c.body.captchaKey = captchaKey
			return c.GetToken()
		}
	}
}

type OAuthTokenResponse struct {
	Token      string `json:"access_token"`
	CaptchaKey string `json:"x_captcha_key"`
	CaptchaUrl string `json:"x_captcha_url"`
	Error      string `json:"error_description"`
}
