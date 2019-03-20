package PhServer

import (
	"log"
	"net/http"
	"strings"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
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

func NewAuthorizeCodeServer(manager oauth2.Manager, rdb *BmRedis.BmRedis) (srv *server.Server) {
	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler(srv, rdb))

	//srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
	//	log.Println("Internal Error:", err.Error())
	//	return
	//})
	//
	//srv.SetResponseErrorHandler(func(re *errors.Response) {
	//	log.Println("Response Error:", re.Error.Error())
	//})

	return
}

func GetInstance(mdb *BmMongodb.BmMongodb, rdb *BmRedis.BmRedis) *server.Server {
	if authServer == nil {
		log.Println("Start ===> AuthorizeCode Server")

		manager := NewAuthorizeCodeManager(mdb, rdb, manage.DefaultAuthorizeCodeTokenCfg)
		authServer = NewAuthorizeCodeServer(manager, rdb)
	}
	return authServer
}

func userAuthorizeHandler(srv *server.Server, rdb *BmRedis.BmRedis) (handler func(w http.ResponseWriter, r *http.Request) (userID string, err error)) {
	handler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		redisDriver := rdb.GetRedisClient()
		defer redisDriver.Close()
		userID, err = redisDriver.Get("LoggedInUserID").Result()
		redisDriver.Del("LoggedInUserID")
		if err != nil || userID == "" {
			//token, ok := srv.BearerAuth(r)
			//if !ok || token == "" {
			redisDriver.Set("ReturnUri", r.Form.Encode(), -1)
			toUrl := strings.Replace(r.URL.Path, "Authorize", "Login", -1)
			w.Header().Set("Location", toUrl)
			w.WriteHeader(http.StatusFound)
		}
		return
	}
	return
}
