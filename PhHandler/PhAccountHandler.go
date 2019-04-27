package PhHandler

import (
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"ph_auth/PhModel"
	"reflect"
	"strings"
	"time"
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
	_ = r.PostForm
	email := r.FormValue("username")
	pwd := r.FormValue("password")

	// Validation Password
	res := PhModel.Account{}
	out := PhModel.Account{}
	cond := bson.M{"email": email, "password": pwd}
	err := h.db.FindOneByCondition(&res, &out, cond)
	if err != nil && out.ID == "" {
		panic("用户名或密码错误")
	}
	if err != nil {
		panic(err.Error())
	}

	redisDriver := h.rd.GetRedisClient()
	defer redisDriver.Close()
	exp := time.Second * 60
	_, err = redisDriver.Set(out.ID+"_login", true, exp).Result()
	if err != nil {
		panic(err.Error())
	}

	a := r.Form
	a.Del("username")
	a.Del("password")

	returnUri := a.Encode()
	toUrl := strings.Replace(r.URL.Path, "AccountValidation", h.Args[1], -1) + "?uid=" + out.ID + "&" + returnUri

	queryForm, _ := url.ParseQuery(r.URL.RawQuery)

	if v := queryForm["status"]; len(v) > 0 && v[0] == "self" {
		fmt.Println(toUrl)
		w.Write([]byte(h.Args[2] + toUrl))
		return 0
	}

	toUrl = strings.Replace(r.URL.Path, "AccountValidation", h.Args[0], -1) + "?uid=" + out.ID + "&" + returnUri

	w.Header().Set("Location", toUrl)
	w.WriteHeader(http.StatusFound)
	return 0
}

func (h PhAccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAccountHandler) GetHandlerMethod() string {
	return h.Method
}
