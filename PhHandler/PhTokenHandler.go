package PhHandler

import (
	"fmt"
	"github.com/PharbersDeveloper/PhAuthServer/PhUnits/array"
	"gopkg.in/oauth2.v3/errors"
	"net/http"
	"reflect"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/mgo.v2/bson"
	"strings"

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
	var scopes []*PhModel.Scope
	cond := bson.M{"_id": bson.ObjectIdHex(token.GetUserID())}
	err = h.db.FindOneByCondition(&res, &out, cond)
	if err != nil {
		panic(err.Error())
	}

	// TODO 人多后有效率问题
	for _, sid := range out.ScopeIDs {
		cond = bson.M{"_id": bson.ObjectIdHex(sid)}
		scopeRes := PhModel.Scope{}
		scopeOut := PhModel.Scope{}
		_ = h.db.FindOneByCondition(&scopeRes, &scopeOut, cond)
		scopes = append(scopes, &scopeOut)
	}

	applyScopes := strings.Split(token.GetScope(), "|") // 申请Scope

	for _, applyScope := range applyScopes {
		detailScope := strings.Split(applyScope, "/")
		level := detailScope[0] // Pharbers 官网 App 单个系统
		action := detailScope[1] // 申请的动作表述

		prefix, applyScopes := scopeSplit(action)
		if len(applyScopes) > 0 {
			truth, _ := singleAppGetScope(scopes, applyScopes, prefix, level)
			if truth { // 输入的Scope超出设置权限，直接跳转到Password登录模式的页面重新验证
				panic(errors.ErrInvalidScope.Error())
			}
		}
	}

	data := map[string]interface{}{
		"expires_in":         int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"refresh_expires_in": token.GetRefreshExpiresIn(),
		"client_id":          token.GetClientID(),
		"user_id":            token.GetUserID(),
		"auth_scope":         token.GetScope(),
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

func singleAppGetScope(accScope []*PhModel.Scope,  applyScopes []string,  prefix, level string) (bool, string) {
	var (
		scope string
		scopeTemp map[string][]string
		temp []string
	)
	scopeTemp = make(map[string][]string)
	for _, v := range accScope {
		temp = append(temp, v.Level + ":" + v.Value)
	}
	for _, applyScope := range applyScopes {
		if array.IsExistItem(prefix + ":" + applyScope, temp) {
			scopeTemp[prefix]= append(scopeTemp[prefix], applyScope)
		} else {
			return true, ""
		}
	}
	body, _ := json.Marshal(scopeTemp)
	scope = fmt.Sprint(level, "/", strings.Trim(strings.ReplaceAll(string(body), `"`, ""), "{}"))
	if scope != "" {
		return false, scope + "|"
	}
	return true, ""

}

func scopeSplit(scope string) (string, []string) {
	var (
		detailScope []string
		scopeSubStr string
	)
	detailScope = strings.Split(scope, ":")
	prefix := detailScope[0]
	if len(detailScope) == 2 {
		scopeSubStr = strings.Trim(detailScope[1], "[]")
		return prefix, strings.Split(scopeSubStr, ",")
	}

	return prefix, nil
}