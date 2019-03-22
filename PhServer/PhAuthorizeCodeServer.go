package PhServer

import (
	"log"
	"net/http"
	"strings"
	"gopkg.in/oauth2.v3"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
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
