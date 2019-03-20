package PhHandler

import (
	"log"
	"net/http"
	"reflect"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/oauth2.v3/server"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"

	"github.com/PharbersDeveloper/PhAuthServer/PhServer"
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
	log.Println("start ===> Generate Token")

	err := h.srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return 0
}

func (h PhTokenHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhTokenHandler) GetHandlerMethod() string {
	return h.Method
}
