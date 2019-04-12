package PhClient

import (
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
)

type PhToken struct {
	oauth2.Token
	Scope     string `json:"scope"`
	AccountID string `json:"account_id"`
}

func ConfigFromURIParameter(r *http.Request) *oauth2.Config {
	queryForm, _ := url.ParseQuery(r.URL.RawQuery)

	config := &oauth2.Config{
		ClientID:     findArrayByKey("client_id", queryForm),
		ClientSecret: findArrayByKey("client_secret", queryForm),
		RedirectURL:  findArrayByKey("redirect_uri", queryForm),
		Scopes:       strings.Split(findArrayByKey("scope", queryForm), "|"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   "http://192.168.100.116:9096/v0/Authorize",
			TokenURL:  "http://192.168.100.116:9096/v0/Token",
		},
	}

	return config
}

func findArrayByKey(key string, values url.Values) string {
	if r := values[key]; len(r) > 0 {
		return r[0]
	}
	return ""
}
