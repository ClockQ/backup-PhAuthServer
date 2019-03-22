package PhHandler

import (
	"net/http"
	"reflect"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/mgo.v2/bson"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"

	"github.com/PharbersDeveloper/PhAuthServer/PhServer"
	"time"
	"encoding/json"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
)

type PhTokenHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
	srv        *server.Server
}

func (h PhTokenHandler) NewTokenHandler(args ...interface{}) PhTokenHandler {
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
		} else {
		}
	}
	sv := PhServer.GetInstance(m, r)

	return PhTokenHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r, srv: sv}
}

func (h PhTokenHandler) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	err := h.srv.HandleTokenRequest(w, r)
	if err != nil {
		panic(err.Error())
	}
	return 0
}

func (h PhTokenHandler) TokenValidation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	token, err := h.srv.ValidationBearerToken(r)
	if err != nil {
		panic(err.Error())
	}

	res := PhModel.Account{}
	out := PhModel.Account{}
	cond := bson.M{"_id": bson.ObjectIdHex(token.GetUserID())}
	err = h.db.FindOneByCondition(&res, &out, cond)
	if err != nil {
		panic(err.Error())
	}

	data := map[string]interface{}{
		"expires_in":         int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"refresh_expires_in": token.GetRefreshExpiresIn(),
		"client_id":          token.GetClientID(),
		"user_id":            token.GetUserID(),
		"auth_scope":         token.GetScope(),
		"all_scope":          out.Scope,
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
	return 0
}

func (h PhTokenHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhTokenHandler) GetHandlerMethod() string {
	return h.Method
}
