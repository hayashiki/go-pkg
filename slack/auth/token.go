package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
	Token interface {
		GetAccessToken(code string) (authResp TokenResponse, err error)
	}

	token struct {
		clientID     string
		clientSecret string
		baseURL      string
		redirectURL  string
		httpClient   *http.Client
	}

	TokenResponse struct {
		Ok          bool       `json:"ok,omitempty"`
		AppID       string     `json:"app_id,omitempty"`
		AuthedUser  AuthedUser `json:"authed_user,omitempty"`
		Scope       string     `json:"scope,omitempty"`
		TokenType   string     `json:"token_type,omitempty"`
		AccessToken string     `json:"access_token,omitempty"`
		BotUserID   string     `json:"bot_user_id,omitempty"`
		Team        TeamInfo   `json:"team,omitempty"`
		Enterprise  string     `json:"enterprise,omitempty"`
		Error       string     `json:"error,omitempty"`
	}

	AuthedUser struct {
		ID string `json:"id,omitempty"`
		Scope string `json:"scope,omitempty"`
		AccessToken string `json:"access_token,omitempty"`
		TokenType string `json:"token_type,omitempty"`
	}

	TeamInfo struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	Option func(*token)
)

func NewToken(clientID, clientSecret, redirectURL string, opts ...Option) Token {
	t := &token{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      SlackDefaultAPIURL,
		redirectURL:  redirectURL,
		httpClient:   http.DefaultClient,
	}

	for _, o := range opts {
		o(t)
	}
	return t
}

func OptionHTTPClient(httpClient *http.Client) Option {
	return func(t *token) {
		t.httpClient = httpClient
	}
}

func (t *token) GetAccessToken(code string) (authResp TokenResponse, err error) {
	v := url.Values{}
	v.Set("client_id", t.clientID)
	v.Set("client_secret", t.clientSecret)
	v.Set("code", code)
	v.Set("redirect_uri", t.redirectURL)

	body := strings.NewReader(v.Encode())

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/oauth.v2.access", t.baseURL), body)
	if err != nil {
		return authResp, fmt.Errorf("error creating slack access token request err=%v", err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", t.clientID, t.clientSecret)))))

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return authResp, fmt.Errorf("failed to get slack authorize response %v", err)
	}

	tokenBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return authResp, fmt.Errorf("failed to get slack access token [%s]: %s", resp.Status, tokenBody)
	}

	err = json.Unmarshal(tokenBody, &authResp)
	if err != nil {
		return authResp, fmt.Errorf("error unmarshal slack Auth response %v", err)
	}

	if !authResp.Ok {
		return authResp, fmt.Errorf(authResp.Error)
	}
	return authResp, nil
}
