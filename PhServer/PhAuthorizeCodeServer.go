package PhServer

import (
	"log"
	"net/http"
	"strings"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/errors"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"gopkg.in/mgo.v2/bson"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"time"
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
	srv.SetUserAuthorizationHandler(userAuthorizeHandler(rdb))
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler(mdb, rdb))

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return
}

func GetInstance(mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis) *server.Server {
	if authServer == nil {
		log.Println("Start ===> AuthorizeCode Server")

		manager := NewAuthorizeCodeManager(mdb, rdb, manage.DefaultAuthorizeCodeTokenCfg)
		authServer = NewAuthorizeCodeServer(manager, mdb, rdb)
	}
	return authServer
}

func userAuthorizeHandler(rdb *BmRedis.BmRedis) (handler func(w http.ResponseWriter, r *http.Request) (userID string, err error)) {
	handler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = r.FormValue("uid")
		result, err := rdb.GetRedisClient().Exists(userID).Result()

		if userID == "" || result == 0 {
			toUrl := strings.Replace(r.URL.Path, "Authorize", "Login", -1)
			returnUri := r.Form.Encode()
			w.Header().Set("Location", toUrl+"?"+returnUri)
			w.WriteHeader(http.StatusFound)
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
		if err == nil && userID!= "" {
			redisDriver := rdb.GetRedisClient()
			defer redisDriver.Close()
			pipe := redisDriver.Pipeline()
			exp := time.Hour * 24 * 3
			pipe.HSet(userID, "nickname", out.Nickname)
			pipe.Expire(userID, exp)
			_, err = pipe.Exec()
		}
		return
	}
	return
}
