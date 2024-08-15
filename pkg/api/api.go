package api

import (
	"net/http"
	"net/http/cookiejar"
)

func init() {
	panic("rest not implemented yet")
}

type Endpoint string

const (
	EndpointUserStatus   Endpoint = "/system/user_status"
	EndpointCheckLogin   Endpoint = "/session/check_login"
	EndpointCheckSession Endpoint = "/session/check_session"
	EndpointCheckSetupOk Endpoint = "/session/check_setup_ok"
)

type Creds struct {
	User string
	Pass string
}

type APIClient struct {
	c     *http.Client
	j     http.CookieJar
	creds *Creds
}

func NewAPIClient(c ...*Creds) *APIClient {
	ac := new(APIClient)
	if len(c) == 1 {
		ac.creds = c[0]
	}
	ac.c = http.DefaultClient
	ac.j, _ = cookiejar.New(nil)
	ac.c.Jar = ac.j
	return ac
}

func (ac *APIClient) WithCookie(cookie *http.Cookie) {
	// ac.j.SetCookies(cookie., []*http.Cookie{cookie})
}

func (ac *APIClient) WithJar(jar http.CookieJar) {
	ac.j = jar
	ac.c.Jar = ac.j
}

func (ac *APIClient) Do(req *http.Request) (*http.Response, error) {
	if ac.creds != nil && req.Header.Get("Authorization") == "" {
		req.SetBasicAuth(ac.creds.User, ac.creds.Pass)
	}
	return ac.c.Do(req)
}

func (ac *APIClient) Valiadate() error {
	panic("unimplemented")
	return nil
}
