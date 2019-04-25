package PhHandler

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/gavv/httpexpect"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/server"
	"net/http"
	"net/http/httptest"
	"net/url"
	"ph_auth/PhServer"
	"ph_auth/PhUnits/yaml"
	"testing"
)

var (
	srv          *server.Server
	tsrv         *httptest.Server
	csrv         *httptest.Server
	clientID     = "5caaf48dd4bc51126652b4c2"
	clientSecret = "5c90db71eeefcc082c0823b2"
	uid          = "5c4552455ee2dd7c36a94a9e"
)

func init() {
	srv = PhServer.GetInstance(mongodbInst(), redisInst())
}

func mongodbInst() *BmMongodb.BmMongodb {
	conf := yaml.LoadConfFromYAML("../resources/resource/service-def.yaml")
	args := conf.Daemons[0].Args
	db := BmMongodb.BmMongodb{
		Host:     args["host"],
		Port:     args["port"],
		Database: args["database"],
	}
	return &db
}

func redisInst() *BmRedis.BmRedis {
	conf := yaml.LoadConfFromYAML("../resources/resource/service-def.yaml")
	args := conf.Daemons[1].Args
	db := BmRedis.BmRedis{
		Host:     args["host"],
		Port:     args["port"],
		Password: args["password"],
		Database: args["database"],
	}
	return &db
}

func clientStore() oauth2.ClientStore {
	clientStore := PhServer.NewAuthorizeCodeClientStore(mongodbInst())
	return clientStore
}

func tokenStore() oauth2.TokenStore {
	tokenStore, _ := PhServer.NewAuthorizeCodeTokenStore(redisInst())
	return tokenStore
}

func testServer(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/authorize":
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			t.Error(err)
		}
	case "/token":
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			t.Error(err)
		}
	}
}

func testClient(t *testing.T, w http.ResponseWriter, r *http.Request, e *httpexpect.Expect) {
	switch r.URL.Path {
	case "/oauth2":
		r.ParseForm()
		code, state := r.Form.Get("code"), r.Form.Get("state")
		if state != "123" {
			t.Error("unrecognized state:", state)
			return
		}
		resObj := e.POST("/token").
			WithFormField("uid", r.Form.Get("uid")).
			WithFormField("redirect_uri", csrv.URL+"/oauth2").
			WithFormField("code", code).
			WithFormField("grant_type", "authorization_code").
			WithFormField("client_id", clientID).
			WithBasicAuth(clientID, clientSecret).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		t.Logf("%#v\n", resObj.Raw())

		validationAccessToken(t, resObj.Value("access_token").String().Raw())
	}
}

func TestRedirectEndPoint(t *testing.T) {

	e := httpexpect.New(t, "http://127.0.0.1:9096/v0")

	e.GET("/ThirdParty").
		WithQuery("client_id", clientID).
		WithQuery("client_secret", clientSecret).
		WithQuery("scope", "ALL").
		WithQuery("state", "xyz").
		WithQuery("redirect_uri", "http://192.168.0.104:4433/oauth-callback").
		Expect().Status(http.StatusOK)

	//outPut := map[string]string{}
	//body,err := ioutil.ReadAll(response.Body())
	//if err != nil {
	//	t.Error(err)
	//}
	//err = json.Unmarshal(body, &outPut)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(outPut)

	//client := &http.Client{}
	//url := fmt.Sprint("http://192.168.0.104:9096/v0/ThirdParty?",
	//						"client_id=5caaf48dd4bc51126652b4c2&",
	//						"client_secret=5c90db71eeefcc082c0823b2&",
	//						"redirect_uri=http://192.168.0.104:4433/oauth-callback")
	//
	//req, err := http.NewRequest("GET", url, nil)
	//
	//response, err := client.Do(req)
	//
	//if err != nil {
	//	t.Error(err)
	//}
	//outPut := map[string]string{}
	//body,err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	t.Error(err)
	//}
	//err = json.Unmarshal(body, &outPut)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(outPut)
}

//func TestTokenEndPoint(t *testing.T) {
//	client := &http.Client{}
//	url := fmt.Sprint("http://192.168.0.104:9096/v0/GenerateAccessToken?",
//		"client_id=5caaf48dd4bc51126652b4c2&",
//		"client_secret=5c90db71eeefcc082c0823b2&",
//		"redirect_uri=http://192.168.0.104:4433/oauth-callback&",
//		"code=50306MWPPVQOOJTONYUNEA")
//
//	req, err := http.NewRequest("GET", url, nil)
//
//
//	response, err := client.Do(req)
//
//	if err != nil {
//		t.Error(err)
//	}
//	outPut := map[string]string{}
//	body,err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		t.Error(err)
//	}
//	err = json.Unmarshal(body, &outPut)
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(outPut)
//}

//func TestAuthorizeCode(t *testing.T) {
//	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		testServer(t, w, r)
//	}))
//	defer tsrv.Close()
//
//	e := httpexpect.New(t, tsrv.URL)
//
//	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		testClient(t, w, r, e)
//	}))
//	defer csrv.Close()
//
//	e.GET("/authorize").
//		WithQuery("uid", uid).
//		WithQuery("response_type", "code").
//		WithQuery("client_id", clientID).
//		WithQuery("scope", "all").
//		WithQuery("state", "123").
//		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
//		Expect().Status(http.StatusOK)
//}

//
//func TestImplicit(t *testing.T) {
//	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		testServer(t, w, r)
//	}))
//	defer tsrv.Close()
//	e := httpexpect.New(t, tsrv.URL)
//
//	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
//	defer csrv.Close()
//
//	manager.MapClientStorage(clientStore(csrv.URL))
//	srv = server.NewDefaultServer(manager)
//	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
//		userID = "000000"
//		return
//	})
//
//	e.GET("/authorize").
//		WithQuery("response_type", "token").
//		WithQuery("client_id", clientID).
//		WithQuery("scope", "all").
//		WithQuery("state", "123").
//		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
//		Expect().Status(http.StatusOK)
//}
//
//func TestPasswordCredentials(t *testing.T) {
//	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		testServer(t, w, r)
//	}))
//	defer tsrv.Close()
//	e := httpexpect.New(t, tsrv.URL)
//
//	manager.MapClientStorage(clientStore(""))
//	srv = server.NewDefaultServer(manager)
//	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
//		if username == "admin" && password == "123456" {
//			userID = "000000"
//			return
//		}
//		err = fmt.Errorf("user not found")
//		return
//	})
//
//	resObj := e.POST("/token").
//		WithFormField("grant_type", "password").
//		WithFormField("username", "admin").
//		WithFormField("password", "123456").
//		WithFormField("scope", "all").
//		WithBasicAuth(clientID, clientSecret).
//		Expect().
//		Status(http.StatusOK).
//		JSON().Object()
//
//	t.Logf("%#v\n", resObj.Raw())
//
//	validationAccessToken(t, resObj.Value("access_token").String().Raw())
//}
//
//func TestClientCredentials(t *testing.T) {
//	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		testServer(t, w, r)
//	}))
//	defer tsrv.Close()
//	e := httpexpect.New(t, tsrv.URL)
//
//	manager.MapClientStorage(clientStore(""))
//
//	srv = server.NewDefaultServer(manager)
//	srv.SetClientInfoHandler(server.ClientFormHandler)
//
//	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
//		t.Log("OAuth 2.0 Error:", err.Error())
//		return
//	})
//
//	srv.SetResponseErrorHandler(func(re *errors.Response) {
//		t.Log("Response Error:", re.Error)
//	})
//
//	srv.SetAllowedGrantType(oauth2.ClientCredentials)
//	srv.SetAllowGetAccessRequest(false)
//	srv.SetExtensionFieldsHandler(func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
//		fieldsValue = map[string]interface{}{
//			"extension": "param",
//		}
//		return
//	})
//	srv.SetAuthorizeScopeHandler(func(w http.ResponseWriter, r *http.Request) (scope string, err error) {
//		return
//	})
//	srv.SetClientScopeHandler(func(clientID, scope string) (allowed bool, err error) {
//		allowed = true
//		return
//	})
//
//	resObj := e.POST("/token").
//		WithFormField("grant_type", "client_credentials").
//		WithFormField("scope", "all").
//		WithFormField("client_id", clientID).
//		WithFormField("client_secret", clientSecret).
//		Expect().
//		Status(http.StatusOK).
//		JSON().Object()
//
//	t.Logf("%#v\n", resObj.Raw())
//
//	validationAccessToken(t, resObj.Value("access_token").String().Raw())
//}

func TestRefreshing(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2":
			r.ParseForm()
			code, state := r.Form.Get("code"), r.Form.Get("state")
			if state != "123" {
				t.Error("unrecognized state:", state)
				return
			}
			jresObj := e.POST("/token").
				WithFormField("redirect_uri", csrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", clientID).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()

			t.Logf("%#v\n", jresObj.Raw())

			validationAccessToken(t, jresObj.Value("access_token").String().Raw())

			resObj := e.POST("/token").
				WithFormField("grant_type", "refresh_token").
				WithFormField("scope", "one").
				WithFormField("refresh_token", jresObj.Value("refresh_token").String().Raw()).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()

			t.Logf("%#v\n", resObj.Raw())

			validationAccessToken(t, resObj.Value("access_token").String().Raw())
		}
	}))
	defer csrv.Close()

	e.GET("/authorize").
		WithQuery("uid", uid).
		WithQuery("response_type", "code").
		WithQuery("client_id", clientID).
		WithQuery("scope", "all").
		WithQuery("state", "123").
		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
		Expect().Status(http.StatusOK)
}

// validation access token
func validationAccessToken(t *testing.T, accessToken string) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	ti, err := srv.ValidationBearerToken(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if ti.GetClientID() != clientID {
		t.Error("invalid access token")
	}
}
