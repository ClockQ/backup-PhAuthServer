package PhClient

import (
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
)
var EndPoint authEndPoint
type PhToken struct {
	oauth2.Token
	Scope     string `json:"scope"`
	AccountID string `json:"account_id"`
}

type authEndPoint struct {
	AuthURL string
	TokenURL string
}

func (e authEndPoint) RegisterEndPoint(value map[string]interface{}) {
	EndPoint = authEndPoint{value["Auth"].(string), value["Token"].(string)}
}

func (e authEndPoint) ConfigFromURIParameter(r *http.Request) *oauth2.Config {
	queryForm, _ := url.ParseQuery(r.URL.RawQuery)

	config := &oauth2.Config{
		ClientID:     findArrayByKey("client_id", queryForm),
		ClientSecret: findArrayByKey("client_secret", queryForm),
		RedirectURL:  findArrayByKey("redirect_uri", queryForm),
		Scopes:       strings.Split(findArrayByKey("scope", queryForm), "|"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  e.AuthURL,
			TokenURL:  e.TokenURL,
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
