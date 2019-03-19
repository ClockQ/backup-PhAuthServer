package PhServer

import (
	"log"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"net/http"
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

func NewAuthorizeCodeServer(manager oauth2.Manager) (srv *server.Server) {
	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

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
		log.Println("Start ===> Generate AuthorizeCode Server")

		manager := NewAuthorizeCodeManager(mdb, rdb, manage.DefaultAuthorizeCodeTokenCfg)
		authServer = NewAuthorizeCodeServer(manager)
	}
	return authServer
}


func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	//w.Header().Set("Location", "/v0/Login")
	//w.WriteHeader(http.StatusFound)

	//toUrl := strings.Replace(r.URL.Path, "AccountValidation", h.Args[0], -1)
	//w.Header().Set("Location", toUrl)
	//store, err := session.Start(nil, w, r)
	//if err != nil {
	//	return
	//}
	//
	//uid, ok := store.Get("LoggedInUserID")
	//if !ok {
	//	if r.Form == nil {
	//		r.ParseForm()
	//	}
	//
	//	store.Set("ReturnUri", r.Form)
	//	store.Save()
	//

	//}
	//
	//userID = uid.(string)
	//store.Delete("LoggedInUserID")
	//store.Save()
	userID = "adbsafd"
	return
}
