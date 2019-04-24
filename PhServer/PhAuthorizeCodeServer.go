package PhServer

import (
	"encoding/json"
	"fmt"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"github.com/PharbersDeveloper/PhAuthServer/PhUnits/array"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"log"
	"net/http"
	"strings"
)

var authServer *server.Server

func NewAuthorizeCodeManager(mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis, authCfg *manage.Config) (manager *manage.Manager) {
	// use authorizeCode Model
	manager = manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(authCfg)

	// token store
	manager.MustTokenStorage(NewAuthorizeCodeTokenStore(rdb))

	// client store
	manager.MapClientStorage(NewAuthorizeCodeClientStore(mdb))
	return
}

func NewAuthorizeCodeServer(manager oauth2.Manager, mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis) (srv *server.Server) {
	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetClientScopeHandler(clientScopeHandler(srv))
	srv.SetUserAuthorizationHandler(userAuthorizeHandler(rdb))
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler(mdb))
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler(mdb))
	return
}

func GetInstance(mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis) *server.Server {
	if authServer == nil {
		manager := NewAuthorizeCodeManager(mdb, rdb, manage.DefaultAuthorizeCodeTokenCfg)
		authServer = NewAuthorizeCodeServer(manager, mdb, rdb)
	}
	return authServer
}

func clientScopeHandler(srv *server.Server) (handler func(clientID, scope string) (allowed bool, err error)) {
	handler = func(clientID, scope string) (allowed bool, err error) {
		_, err = srv.Manager.GetClient(clientID)
		if err != nil {
			return
		}
		allowed = true
		return
	}
	return
}

func userAuthorizeHandler(rdb *BmRedis.BmRedis) (handler func(w http.ResponseWriter, r *http.Request) (userID string, err error)) {
	handler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		//queryForm, _ := url.ParseQuery(r.URL.RawQuery)
		userID = r.FormValue("uid")
		redisDriver := rdb.GetRedisClient()
		result, err := redisDriver.Exists(userID + "_login").Result()
		redisDriver.Del(userID)
		if userID == "" || result == 0 {
			log.Println("用户未登录或操作超时，转至登录页")
			userID = ""
			toUrl := strings.Replace(r.URL.Path, "Authorize", "Login", -1)
			returnUri := r.Form.Encode()
			w.Header().Set("Location", toUrl+"?"+returnUri)
			w.WriteHeader(http.StatusFound)
			return
		}
		return
	}
	return
}

func passwordAuthorizationHandler(mdb *BmMongodb.BmMongodb) (handler func(username, password string) (userID string, err error)) {
	handler = func(email, pwd string) (userID string, err error) {
		res := PhModel.Account{}
		out := PhModel.Account{}
		cond := bson.M{"email": email, "password": pwd}
		_ = mdb.FindOneByCondition(&res, &out, cond)

		userID = out.ID
		if userID == "" {
			log.Println("用户使用密码验证，但登录失败")
		}

		return
	}
	return
}

func authorizeScopeHandler(mdb *BmMongodb.BmMongodb) (handler func(w http.ResponseWriter, r *http.Request) (scope string, err error)) {
	handler = func(w http.ResponseWriter, r *http.Request) (scope string, err error) {
		var scopes []*PhModel.Scope
		accRes := PhModel.Account{}
		accOut := PhModel.Account{}
		userID := r.FormValue("uid")
		applyScopes := strings.Split(r.FormValue("scope"), " ") // 申请Scope
		cond := bson.M{"_id": bson.ObjectIdHex(userID)}
		err = mdb.FindOneByCondition(&accRes, &accOut, cond)

		if err != nil {
			w.Header().Set("Location", "http://www.baidu.com") // TODO：暂定跳转百度
			w.WriteHeader(http.StatusFound)
			return
		}

		// TODO 人多后有效率问题
		for _, sid := range accOut.ScopeIDs {
			cond = bson.M{"_id": bson.ObjectIdHex(sid)}
			scopeRes := PhModel.Scope{}
			scopeOut := PhModel.Scope{}
			_ = mdb.FindOneByCondition(&scopeRes, &scopeOut, cond)
			scopes = append(scopes, &scopeOut)
		}

		for _, applyScope := range applyScopes {
			detailScope := strings.Split(applyScope, "/")
			level := detailScope[0] // Pharbers 官网 App 单个系统
			action := detailScope[1] // 申请的动作表述

			prefix, applyScopes := scopeSplit(action)
			if len(applyScopes) > 0 {
				truth, result := singleAppGetScope(scopes, applyScopes, prefix, level)
				if truth { // App登入时输入的Scope超出设置权限，直接跳转到Password登录模式的页面重新验证
					w.Header().Set("Location", "http://www.baidu.com") // TODO：暂定跳转百度
					w.WriteHeader(http.StatusFound)
				} else {
					scope += result
				}
			} else {
				scope += topLevelGetScope(scopes, level, prefix)
			}
		}
		scope = scope[:strings.LastIndex(scope, "|")]
		return
	}
	return
}

func topLevelGetScope(accScope []*PhModel.Scope, level,prefix string) string {
		var (
			scope string
			scopeTemp map[string][]string
		)
	scopeTemp = make(map[string][]string)

		for _, s := range accScope {
			if s.Level == prefix {
				scopeTemp[prefix]= append(scopeTemp[prefix], s.Value)
			}
		}

		body, _ := json.Marshal(scopeTemp)
		scope = fmt.Sprint(level, "/", strings.Trim(strings.ReplaceAll(string(body), `"`, ""), "{}"), "|")
		return scope
}

func singleAppGetScope(accScope []*PhModel.Scope,  applyScopes []string,  prefix, level string) (bool, string) {
	var (
		scope string
		scopeTemp map[string][]string
		temp []string
	)
	scopeTemp = make(map[string][]string)
	for _, v := range accScope {
		temp = append(temp, v.Level + ":" + v.Value)
	}
	for _, applyScope := range applyScopes {
		if array.IsExistItem(prefix + ":" + applyScope, temp) {
			scopeTemp[prefix]= append(scopeTemp[prefix], applyScope)
		} else {
			return true, ""
		}
	}
	body, _ := json.Marshal(scopeTemp)
	scope = fmt.Sprint(level, "/", strings.Trim(strings.ReplaceAll(string(body), `"`, ""), "{}"))
	if scope != "" {
		return false, scope + "|"
	}
	return true, ""

}

func scopeSplit(scope string) (string, []string) {
	var (
		detailScope []string
		scopeSubStr string
	)
	detailScope = strings.Split(scope, ":")
	prefix := detailScope[0]
	if len(detailScope) == 2 {
		scopeSubStr = strings.Trim(detailScope[1], "[]")
		return prefix, strings.Split(scopeSubStr, ",")
	}

	return prefix, nil
}