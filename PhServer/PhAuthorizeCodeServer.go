package PhServer

import (
	"encoding/json"
	"fmt"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
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
		accRes := PhModel.Account{}
		accOut := PhModel.Account{}
		userID := r.FormValue("uid")
		applyScopes := strings.Split(r.FormValue("scope"), " ") // 申请Scope
		cond := bson.M{"_id": bson.ObjectIdHex(userID)}
		_ = mdb.FindOneByCondition(&accRes, &accOut, cond)

		for _, applyScope := range applyScopes {
			detailScope := strings.Split(applyScope, "/")
			level := detailScope[0] // Pharbers 官网 App 单个系统
			action := detailScope[1] // 申请的动作表述
			accScopes := strings.Split(accOut.Scope, "|")

			prefix, scopes := scopeSplit(action)

			if len(scopes) > 0 {
				truth, result := singleAppGetScope(accScopes, scopes, prefix, level)
				if truth { // App登入时输入的Scope超出设置权限，直接跳转到Password登录模式的页面重新验证
					w.Header().Set("Location", "http://www.baidu.com") // TODO：暂定跳转百度
					w.WriteHeader(http.StatusFound)
				} else {
					scope += result
				}
			} else {
				scope += topLevelGetScope(accScopes, action, level)
			}
		}
		return
	}
	return
}

func topLevelGetScope(accScope []string, applyScope, level string) string {
		var scope string
		for _, v := range accScope {
			if strings.Contains(v, applyScope) {
				scope += fmt.Sprint(level, "/", v, "|")
			}
		}
		return scope
}

func singleAppGetScope(accScope, scopes []string,  prefix, level string) (bool, string) {
	var (
		scope string
		temp map[string][]string
	)
	temp = make(map[string][]string)
	for _, applyScope := range scopes {
		for _, v := range accScope {
			if prefix == strings.Split(v, ":")[0] {
				if strings.Contains(v, applyScope) {
					key := strings.Split(v, ":")[0]
					temp[key]= append(temp[key], applyScope)
				} else {
					return true, ""
				}
			}
		}
	}
	body, _ := json.Marshal(temp)
	scope = strings.Trim(strings.ReplaceAll(string(body), `"`, ""), "{}")
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