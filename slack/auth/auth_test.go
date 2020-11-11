package auth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type getTokenTransport struct {
	StatusCode int
	Body       string
}

func (g getTokenTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: g.StatusCode,
		Body:       ioutil.NopCloser(strings.NewReader(g.Body)),
	}
	return resp, nil
}

func TestNewAuth(t *testing.T) {
	type args struct {
		clientID    string
		redirectURL string
		state       string
		scopes      []string
	}
	tests := []struct {
		name string
		args args
		want Auth
	}{
		{
			name: "noError",
			args: args{
				clientID:    "1",
				redirectURL: "http://localhost/callback",
				state:       "10b2d11b-7064-4bfe-92e1-a1cd96d77713",
				scopes:      []string{"chat:write"},
			},
			want: &auth{
				clientID:    "1",
				redirectURL: "http://localhost/callback",
				baseURL:     "https://slack.com",
				state:       "10b2d11b-7064-4bfe-92e1-a1cd96d77713",
				scopes:      []string{"chat:write"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuth(tt.args.clientID, tt.args.redirectURL, tt.args.state, tt.args.scopes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want *AuthorizeResponse
	}{
		{
			name: "noError",
			// http://localhost:8080/slack/callback?code=12345678.12345678.1234567890pooiuytrewq&state=10b2d11b-7064-4bfe-92e1-a1cd96d77713
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/slack/callback?code=12345678.12345678.1234567890pooiuytrewq&state=10b2d11b-7064-4bfe-92e1-a1cd96d77713", nil),
			},
			want: &AuthorizeResponse{
				Code:  "12345678.12345678.1234567890pooiuytrewq",
				State: "10b2d11b-7064-4bfe-92e1-a1cd96d77713",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRequest(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Redirect(t *testing.T) {
	type fields struct {
		State       string
		clientID    string
		redirectURL string
		baseURL     string
		scopes      []string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "noError",
			fields: fields{
				State:       "hoge",
				clientID:    "1",
				redirectURL: "http://localhost/callback",
				baseURL:     "https://slack.com",
				scopes:      []string{"chat:write", "users:read"},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewAuth(tt.fields.clientID, tt.fields.redirectURL, tt.fields.State, tt.fields.scopes)
			if err := c.Redirect(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Redirect() error = %v, wantErr %v", err, tt.wantErr)
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			if status := w.Code; status != http.StatusTemporaryRedirect {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusTemporaryRedirect)
			}

			expectURL := "https://slack.com/oauth/v2/authorize?client_id=1&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&scope=chat%3Awrite%2Cusers%3Aread&state=hoge"
			if actualURL := tt.args.w.Header().Get("location"); actualURL != expectURL {
				t.Errorf("handler returned wrong status code: got %v want %v",
					actualURL, expectURL)
			}
		})
	}
}

func Test_auth_UserRedirect(t *testing.T) {
	type fields struct {
		State       string
		clientID    string
		redirectURL string
		baseURL     string
		userScopes  []string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "noError",
			fields: fields{
				State:       "hoge",
				clientID:    "1",
				redirectURL: "http://localhost/callback",
				baseURL:     "https://slack.com",
				userScopes:  []string{"user.identify"},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewUserAuth(tt.fields.clientID, tt.fields.redirectURL, tt.fields.State, tt.fields.userScopes)
			if err := c.Redirect(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Redirect() error = %v, wantErr %v", err, tt.wantErr)
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			if status := w.Code; status != http.StatusTemporaryRedirect {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusTemporaryRedirect)
			}

			expectURL := "https://slack.com/oauth/v2/authorize?client_id=1&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&state=hoge&user_scope=user.identify"
			if actualURL := tt.args.w.Header().Get("location"); actualURL != expectURL {
				t.Errorf("handler returned wrong status code: got %v want %v",
					actualURL, expectURL)
			}
		})
	}
}

func Test_token_GetAccessToken(t1 *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		baseURL      string
		redirectURL  string
		httpClient   *http.Client
	}
	type args struct {
		code string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAuthResp TokenResponse
		wantErr      bool
	}{
		{
			name: "noError",
			fields: fields{
				clientID:     "1",
				clientSecret: "1",
				redirectURL:  "http://localhost",
			},
			args: args{
				code: "code",
			},
			wantAuthResp: TokenResponse{
				Ok:    true,
				AppID: "A015N9ABCD",
				AuthedUser: AuthedUser{
					ID: "UD7AKABCD",
				},
				Scope:       "chat:write,users:read",
				TokenType:   "bot",
				AccessToken: "yoxb-dummy-ckoHABCD4VOWr",
				BotUserID:   "U016FUXABCD",
				Team: TeamInfo{
					ID:   "T0CB2AABCD",
					Name: "dummy",
				},
				Enterprise: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			httpClient := &http.Client{}
			httpClient.Transport = &getTokenTransport{
				StatusCode: http.StatusOK,
				Body:       `{"ok":true,"app_id":"A015N9ABCD","authed_user":{"id":"UD7AKABCD"},"scope":"chat:write,users:read","token_type":"bot","access_token":"yoxb-dummy-ckoHABCD4VOWr","bot_user_id":"U016FUXABCD","team":{"id":"T0CB2AABCD","name":"dummy"},"enterprise":null}`,
			}
			t := NewToken(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURL, OptionHTTPClient(httpClient))
			gotAuthResp, err := t.GetAccessToken(tt.args.code)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAuthResp, tt.wantAuthResp) {
				t1.Errorf("GetAccessToken() gotAuthResp = %v, want %v", gotAuthResp, tt.wantAuthResp)
			}
		})
	}
}
