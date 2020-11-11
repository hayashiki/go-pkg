package auth

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	SlackDefaultAPIURL = "https://slack.com"
)

type (
	Auth interface {
		Redirect(w http.ResponseWriter, r *http.Request) error
		CallbackHTML(TeamID string) string
	}

	auth struct {
		state       string
		clientID    string
		redirectURL string
		baseURL     string
		scopes      []string
		userScopes  []string
	}

	AuthorizeResponse struct {
		Code  string
		State string
	}
)

func (c *auth) accessURL() (string, error) {
	u, err := url.Parse(c.baseURL)
	u.Path = "oauth/v2/authorize"
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Add("client_id", c.clientID)
	v.Add("redirect_uri", c.redirectURL)
	v.Add("state", c.state)
	if len(c.scopes) != 0 {
		v.Add("scope", strings.Join(c.scopes, ","))
	}
	if len(c.userScopes) != 0 {
		v.Add("user_scope", strings.Join(c.userScopes, ","))
	}
	u.RawQuery = v.Encode()

	return u.String(), nil
}

// Redirect redirect to request url
func (c *auth) Redirect(w http.ResponseWriter, r *http.Request) error {
	urlStr, err := c.accessURL()
	if err != nil {
		return err
	}
	http.Redirect(w, r, urlStr, http.StatusTemporaryRedirect)
	return nil
}

func ParseRequest(r *http.Request) *AuthorizeResponse {
	return &AuthorizeResponse{
		Code:  r.FormValue("code"),
		State: r.FormValue("state"),
	}
}

func (c *auth) CallbackHTML(TeamID string) string {
	panic("implement me")
}

func NewAuth(clientID, redirectURL, state string, scopes []string) Auth {
	return &auth{
		clientID:    clientID,
		redirectURL: redirectURL,
		baseURL:     SlackDefaultAPIURL,
		scopes:      scopes,
		state:       state,
	}
}

func NewUserAuth(clientID, redirectURL, state string, userScopes []string) Auth {
	return &auth{
		clientID:    clientID,
		redirectURL: redirectURL,
		baseURL:     SlackDefaultAPIURL,
		userScopes:  userScopes,
		state:       state,
	}
}
