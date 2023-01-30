package yandex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	fmt.Println(string(body))
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	if len(response.Error) > 0 {
		err = errors.New(response.Error)
	}
	return
}

type OAuthTokenResponse struct {
	Token string `json:"access_token"`
	Error string `json:"error"`
}
