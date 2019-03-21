package PhServer

import (
		"net/http"
	"strings"
	"gopkg.in/oauth2.v3"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/errors"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"github.com/PharbersDeveloper/PhAuthServer/PhUnits/array"
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
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler(mdb))
	srv.SetUserAuthorizationHandler(userAuthorizeHandler(rdb))
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler(mdb, rdb))

	return
}

func GetInstance(mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis) *server.Server {
	if authServer == nil {
		manager := NewAuthorizeCodeManager(mdb, rdb, manage.DefaultAuthorizeCodeTokenCfg)
		authServer = NewAuthorizeCodeServer(manager, mdb, rdb)
	}
	return authServer
}

func userAuthorizeHandler(rdb *BmRedis.BmRedis) (handler func(w http.ResponseWriter, r *http.Request) (userID string, err error)) {
	handler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = r.FormValue("uid")
		redisDriver := rdb.GetRedisClient()
		result, err := redisDriver.Exists(userID).Result()
		if userID == "" || result == 0 {
			toUrl := strings.Replace(r.URL.Path, "Authorize", "Login", -1)
			returnUri := r.Form.Encode()
			w.Header().Set("Location", toUrl+"?"+returnUri)
			w.WriteHeader(http.StatusFound)
		}
		redisDriver.Del(userID)
		return
	}
	return
}

func authorizeScopeHandler(mdb *BmMongodb.BmMongodb) (handler func(w http.ResponseWriter, r *http.Request) (scope string, err error)) {
	handler = func(w http.ResponseWriter, r *http.Request) (scope string, err error) {
		// Validation Scope
		uid := r.FormValue("uid")
		res := PhModel.Account{}
		out := PhModel.Account{}
		cond := bson.M{"_id": bson.ObjectIdHex(uid)}
		err = mdb.FindOneByCondition(&res, &out, cond)

		bl := false
		if array.IsExistItem("ALL", strings.Split(out.Scope, "#")) {
			bl = true
		} else {
			bl = array.IsExistItem(scope, strings.Split(out.Scope, "#"))
		}

		if bl == false {
			err = errors.ErrInvalidScope
		}

		return
	}
	return
}

func passwordAuthorizationHandler(mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis) (handler func(username, password string) (userID string, err error)) {
	handler = func(email, pwd string) (userID string, err error) {
		res := PhModel.Account{}
		out := PhModel.Account{}
		cond := bson.M{"email": email, "password": pwd}
		err = mdb.FindOneByCondition(&res, &out, cond)

		userID = out.ID
		return
	}
	return
}
