package PhServer

import (
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/manyminds/api2go"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"log"
	"net/http"
	"ph_auth/PhModel"
	"strings"
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
	srv.SetClientScopeHandler(clientScopeHandler(srv))
	srv.SetUserAuthorizationHandler(userAuthorizeHandler(rdb))
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler(mdb))
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler(mdb))
	srv.SetAccessTokenExpHandler(accessTokenExpHandler())
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
		accRes := PhModel.Account{}
		accOut := PhModel.Account{}
		userID := r.FormValue("uid")
		applyScope := r.FormValue("scope") // 申请Scope
		cond := bson.M{"_id": bson.ObjectIdHex(userID)}
		err = mdb.FindOneByCondition(&accRes, &accOut, cond)
		if err != nil {
			w.Header().Set("Location", "http://www.pharbers.com") // TODO：暂定跳转官网
			w.WriteHeader(http.StatusFound)
			return
		}

		empModel := PhModel.Employee{}
		err = mdb.FindOneByCondition(&empModel, &empModel, bson.M{"_id": bson.ObjectIdHex(accOut.EmployeeID)})
		if err != nil {
			w.Header().Set("Location", "http://www.pharbers.com") // TODO：暂定跳转官网
			w.WriteHeader(http.StatusFound)
			return
		}

		scopeReq := api2go.Request{
			QueryParams: map[string][]string{
				"group-id": {empModel.GroupID},
			},
		}
		scopeIn := PhModel.Scope{}
		var scopeModels []PhModel.Scope
		err = mdb.FindMulti(scopeReq, &scopeIn, &scopeModels, -1, -1)
		if err == nil {
			for i, iter := range scopeModels {
				mdb.ResetIdWithId_(&iter)
				scopeModels[i] = iter
			}
		} else {
			w.Header().Set("Location", "http://www.pharbers.com") // TODO：暂定跳转官网
			w.WriteHeader(http.StatusFound)
			return
		}

		detailScope := strings.Split(applyScope, "/")
		level := detailScope[0]       // Pharbers 官网 App 单个系统
		applyAccess := detailScope[1] // 申请的动作表述
		accessed, scop := checkAccessScope(applyAccess, scopeModels)	//允许访问(判断是否有授权)
		indate := checkIndateScope(applyAccess, scopeModels)	//有效期内(判断授权是否过期)
		if accessed && indate {
			scope = level + "/" + scop
			return
		} else {
			w.Header().Set("Location", "http://www.pharbers.com") // TODO：暂定跳转官网
			w.WriteHeader(http.StatusFound)
			return
		}

	}
	return
}

// 设置AccessToken过期时间
func accessTokenExpHandler () (handler func(w http.ResponseWriter, r *http.Request)(exp time.Duration, err error))  {
	handler = func(w http.ResponseWriter, r *http.Request) (exp time.Duration, err error) {
		h := time.Hour * 8
		fmt.Println(h.Minutes())
		return h, nil
	}
	return handler
}

func checkAccessScope(applyAccess string, accScopes []PhModel.Scope) (accessed bool, scope string) {

	accessed = false
	for _, v := range accScopes {
		if v.Access == applyAccess {
			accessed = true
			scope = fmt.Sprint(v.Access, ":", v.Operation, "#", v.Expired)
			return
		}
	}
	return
}

func checkIndateScope(applyAccess string, accScopes []PhModel.Scope) (indate bool) {

	indate = false
	for _, v := range accScopes {
		if v.Access == applyAccess {
			now := float64(time.Now().UnixNano() / 1e6)
			if now < v.Expired {
				indate = true
				return
			}
		}
	}
	return
}
