package PhHandler

import (
	"encoding/json"
	"github.com/PharbersDeveloper/PhAuthServer/PhClient"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
)

type PhUserAgentHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h PhUserAgentHandler) NewUserAgentHandle(args ...interface{}) PhUserAgentHandler {
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

	return PhUserAgentHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h PhUserAgentHandler) ThirdParty(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {

	config := PhClient.ConfigFromURIParameter(r)
	url := config.AuthCodeURL("xyz")
	response := map[string]interface{} {
		"redirect-uri": url,
	}
	// 这个地方再想想
	enc := json.NewEncoder(w)
	enc.Encode(response)
	//http.Redirect(w, r, url, http.StatusFound)
	return 0
}

func (h PhUserAgentHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhUserAgentHandler) GetHandlerMethod() string {
	return h.Method
}
