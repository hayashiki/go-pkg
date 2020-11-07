package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hayashiki/go-pkg/slack/auth"
	"net/http"
	"os"
	"time"
)

var (
	ClientID = os.Getenv("SLACK_CLIENT_ID")
	ClientSecret = os.Getenv("SLACK_SECRET_ID")
	CallbackURL = os.Getenv("SLACK_REDIRECT_URL")
)

func main() {
	// open http://localhost:8080/auth
	http.HandleFunc("/slack/auth", Authorize)
	http.HandleFunc("/slack/callback", Callback)

	http.ListenAndServe(":8080", nil)
}

func Authorize(w http.ResponseWriter, r *http.Request)  {
	state := uuid.New().String()
	scopes := []string{"chat:write", "users:read", "users.profile:read", "channels:read", "groups:read", "im:read", "mpim:read", "commands"}
	c := auth.NewAuth(ClientID, CallbackURL, state, scopes)

	http.SetCookie(w, &http.Cookie{
		Name:       "state",
		Value:      state,
		Expires:    time.Now().Add(60 * time.Second),
	})

	// https://slack.com/oauth/v2/authorize?client_id=123456789.14712345678&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fslack%2Fcallback&scope=chat%3Awrite%2Cusers%3Aread%2Cusers.profile%3Aread%2Cchannels%3Aread%2Cgroups%3Aread%2Cim%3Aread%2Cmpim%3Aread%2Ccommands&state=1be65a10-1f6b-4013-8c85-03f728d5ba7d
	c.Redirect(w, r)
}

// eg. http://localhost:8080/slack/callback?code=123456789.14712345678.1234567890qwertyuioo&state=10b2d11b-7064-4bfe-92e1-a1cd96d77713
func Callback(w http.ResponseWriter, r *http.Request) {
	resp := auth.ParseRequest(r)

	state, err := r.Cookie("state")
	if err != nil {
		fmt.Fprintf(w, "err=%v", err)
		return
	}
	if resp.State != state.Value {
		fmt.Fprintf(w, "err=%v", err)
		return
	}

	c := auth.NewToken(ClientID, ClientSecret, CallbackURL)
	accessToken, err := c.GetAccessToken(resp.Code)
	if err != nil {
		fmt.Fprintf(w, "err=%v", err)
		return
	}
	fmt.Fprintf(w, "token:%v", accessToken)
}
