package PhHandler

import (
	"net/http"
	"reflect"
	"encoding/json"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"github.com/julienschmidt/httprouter"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"log"
	)

type PhAccountHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h PhAccountHandler) NewAccountHandler(args ...interface{}) PhAccountHandler {
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

	return PhAccountHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h PhAccountHandler) AccountValidation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	log.Println("Start ===> Account Validation")

	_ = r.PostForm
	email := r.FormValue("username")
	pwd := r.FormValue("password")

	res := PhModel.Account{}
	out := PhModel.Account{}
	email = "zyqi@pharbers.com"
	pwd = "4297f44b13955235245b2497399d7a93"
	cond := bson.M{"email": email, "password": pwd}
	err := h.db.FindOneByCondition(&res, &out, cond)

	if err == nil && out.ID != "" {
		//redisDriver := h.rd.GetRedisClient()
		//defer redisDriver.Close()
		//pipe := redisDriver.Pipeline()
		//exp := time.Hour * 24 * 3
		//pipe.HSet(out.ID, "nickname", out.Nickname)
		//pipe.Expire(out.ID, exp)
		//_, err = pipe.Exec()

		toUrl := strings.Replace(r.URL.Path, "AccountValidation", h.Args[0], -1)
		a := r.Form
		a.Del("username")
		a.Del("password")
		returnUri := a.Encode()
		w.Header().Set("Location", toUrl+"?uid="+out.ID+"&"+returnUri)
		w.WriteHeader(http.StatusFound)
		return 0
	} else {
		response := map[string]interface{}{
			"status": "error",
			"result": nil,
			"error":  "账户或密码错误！",
		}
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}
}

func (h PhAccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAccountHandler) GetHandlerMethod() string {
	return h.Method
}
