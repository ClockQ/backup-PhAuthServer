package PhHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3/server"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/PharbersDeveloper/PhAuthServer/PhClient"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"github.com/PharbersDeveloper/PhAuthServer/PhServer"
	"reflect"
	"time"
)

type PhAuthorizeHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
	srv        *server.Server
}

func (h PhAuthorizeHandler) NewAuthorizeHandler(args ...interface{}) PhAuthorizeHandler {
	var m *BmMongodb.BmMongodb
	var r *BmRedis.BmRedis
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				}
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		}
	}
	sv := PhServer.GetInstance(m, r)

	return PhAuthorizeHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r, srv: sv}
}

func (h PhAuthorizeHandler) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	err := h.srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		panic(err.Error())
	}
	return 0
}

func (h PhAuthorizeHandler) GenerateAccessToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	queryParameter, _ := url.ParseQuery(r.URL.RawQuery)

	config := PhClient.EndPoint.ConfigFromURIParameter(r)

	accessToken, err := config.Exchange(context.Background(), queryParameter["code"][0])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}

	// 获取AuthServer 存入的UserID
	tokenUid, _ := h.RdGetValueByKey(accessToken.AccessToken)
	initialToken, _ := h.RdGetValueByKey(tokenUid)
	var oauthPrototype map[string]interface{}
	json.Unmarshal([]byte(initialToken), &oauthPrototype)

	phToken := PhClient.PhToken{
		Scope:     accessToken.Extra("scope").(string),
		AccountID: oauthPrototype["UserID"].(string),
	}
	phToken.AccessToken = accessToken.AccessToken
	phToken.RefreshToken = accessToken.RefreshToken
	phToken.Expiry = accessToken.Expiry
	phToken.TokenType = accessToken.TokenType

	// 存入Redis RefreshToken
	//err = h.PushValueByKey("RefreshToken_"+phToken.RefreshToken, &phToken)
	err = h.PushValueByKeyAndExpire(fmt.Sprint("RefreshToken_", phToken.RefreshToken), &phToken, time.Hour * 12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	err = e.Encode(phToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}

	return 0
}

func (h PhAuthorizeHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	queryForm, _ := url.ParseQuery(r.URL.RawQuery)

	refreshToken := queryForm["refresh_token"][0]
	if len(refreshToken) <= 0 {
		http.Error(w, "refresh_token invalid", http.StatusBadRequest)
		return 1
	}

	config := PhClient.EndPoint.ConfigFromURIParameter(r)

	token := &oauth2.Token{}
	tokenResult, err := h.RdGetValueByKey("RefreshToken_" + refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}
	json.Unmarshal([]byte(tokenResult), &token)

	token.Expiry = time.Now()

	accessToken, err := config.TokenSource(context.Background(), token).Token()

	tokenUid, _ := h.RdGetValueByKey(accessToken.AccessToken)
	initialToken, _ := h.RdGetValueByKey(tokenUid)

	var oauthPrototype map[string]interface{}
	json.Unmarshal([]byte(initialToken), &oauthPrototype)

	phToken := PhClient.PhToken{
		Scope:     accessToken.Extra("scope").(string),
		AccountID: oauthPrototype["UserID"].(string),
	}
	phToken.AccessToken = accessToken.AccessToken
	phToken.RefreshToken = accessToken.RefreshToken
	phToken.Expiry = accessToken.Expiry
	phToken.TokenType = accessToken.TokenType

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}

	// 存入Redis RefreshToken
	//err = h.PushValueByKey("RefreshToken_"+phToken.RefreshToken, &phToken)
	err = h.PushValueByKeyAndExpire(fmt.Sprint("RefreshToken_", phToken.RefreshToken), &phToken, time.Hour * 12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}

	defer h.RdDeleteToken("RefreshToken_" + refreshToken)

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(phToken)
	return 0
}

func (h PhAuthorizeHandler) PasswordLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	body, err := ioutil.ReadAll(r.Body)

	var parameter map[string]interface{}
	json.Unmarshal(body, &parameter)

	config := PhClient.EndPoint.ConfigFromURIParameter(r)
	token, err := config.PasswordCredentialsToken(context.Background(), parameter["username"].(string), parameter["password"].(string))

	// 获取AuthServer 存入的UserID
	tokenUUID, _ := h.RdGetValueByKey(token.AccessToken)
	initialToken, _ := h.RdGetValueByKey(tokenUUID)
	var oauthPrototype map[string]interface{}
	json.Unmarshal([]byte(initialToken), &oauthPrototype)

	accRes := PhModel.Account{}
	accOut := PhModel.Account{}
	cond := bson.M{"_id": bson.ObjectIdHex(oauthPrototype["UserID"].(string))}
	err = h.db.FindOneByCondition(&accRes, &accOut, cond)
	if err != nil {
		panic(err.Error())
	}

	empModel := PhModel.Employee{}
	err = h.db.FindOneByCondition(&empModel, &empModel, bson.M{"_id": bson.ObjectIdHex(accOut.EmployeeID)})
	if err != nil {
		panic(err.Error())
	}

	scopeReq := api2go.Request{
		QueryParams: map[string][]string{
			"group-id": {empModel.GroupID},
		},
	}
	scopeIn := PhModel.Scope{}
	var scopeModels []PhModel.Scope
	err = h.db.FindMulti(scopeReq, &scopeIn, &scopeModels, -1, -1)
	if err == nil {
		for i, iter := range scopeModels {
			h.db.ResetIdWithId_(&iter)
			scopeModels[i] = iter
		}
	} else {
		panic(err.Error())
	}

	phToken := PhClient.PhToken{
		Scope:     token.Extra("scope").(string),
		AccountID: oauthPrototype["UserID"].(string),
	}

	if len(scopeModels) != 0 {
		phToken.Scope = makeScopeStr(scopeModels)

		if len(oauthPrototype) == 0 {
			panic("oauthPrototype is empty")
		}
		oauthPrototype["Scope"] = phToken.Scope
		expired := int64(oauthPrototype["AccessExpiresIn"].(float64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 1
		}

		h.PushValueByKeyAndExpire(tokenUUID, oauthPrototype, 10 * time.Duration(expired))

	}

	phToken.AccessToken = token.AccessToken
	phToken.RefreshToken = token.RefreshToken
	phToken.Expiry = token.Expiry
	phToken.TokenType = token.TokenType

	// 存入Redis RefreshToken
	//err = h.PushValueByKey("RefreshToken_"+phToken.RefreshToken, &phToken)
	err = h.PushValueByKeyAndExpire(fmt.Sprint("RefreshToken_", phToken.RefreshToken), &phToken, time.Hour * 12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	err = e.Encode(phToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 1
	}
	return 0
}

func (h PhAuthorizeHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAuthorizeHandler) GetHandlerMethod() string {
	return h.Method
}

func (h PhAuthorizeHandler) PushValueByKeyAndExpire(key string, value interface{}, expiration time.Duration) error {
	jsonToken, _ := json.Marshal(value)

	client := h.rd.GetRedisClient()
	defer client.Close()

	pipe := client.Pipeline()
	pipe.Set(key, string(jsonToken), expiration)

	_, err := pipe.Exec()
	return err
}

func (h PhAuthorizeHandler) RdGetValueByKey(key string) (string, error) {
	client := h.rd.GetRedisClient()
	defer client.Close()

	result, err := client.Get(key).Result()

	if err != nil {
		return "", err
	}
	return result, nil
}

func (h PhAuthorizeHandler) RdDeleteToken(key string) {
	client := h.rd.GetRedisClient()
	defer client.Close()

	pipe := client.Pipeline()

	pipe.Del(key)

	pipe.Exec()
}

func makeScopeStr(scopeArr []PhModel.Scope) (scope string) {
	if len(scopeArr) == 0 {
		panic("scopeArr is empty")
	}
	scope = fmt.Sprint("APP/", scopeArr[0].Access, ":", scopeArr[0].Operation, "#", scopeArr[0].Expired)

	for _, v := range scopeArr[1:] {
		scope = fmt.Sprint(scope, ",", v.Access, ":", v.Operation, "#", v.Expired)
	}

	return
}
