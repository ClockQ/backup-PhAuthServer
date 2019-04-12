package PhHandler

import (
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"net/http"
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

	// Validation Scope
	// 没啥用
	//scope := r.FormValue("scope")
	//bl := false
	//if array.IsExistItem("ALL", strings.Split(out.Scope, "|")) {
	//	bl = true
	//} else {
	//	for _, v := range strings.Split(scope, "|") {
	//		if array.IsExistItem(v, strings.Split(out.Scope, "|")) {
	//			bl = true
	//			break
	//		}
	//	}
	//}
	//
	//if bl == false {
	//	panic(fmt.Sprintf("登录失败, 传入 'scope = %s' 错误，或用户没有 '%s' 的权限", scope, scope))
	//}

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
	toUrl := strings.Replace(r.URL.Path, "AccountValidation", h.Args[0], -1)
	w.Header().Set("Location", toUrl+"?uid="+out.ID+"&"+returnUri)
	w.WriteHeader(http.StatusFound)
	return 0
}

func (h PhAccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAccountHandler) GetHandlerMethod() string {
	return h.Method
}
