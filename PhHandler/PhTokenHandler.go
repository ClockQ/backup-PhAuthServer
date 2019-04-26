package PhHandler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/server"
	"net/http"
	"ph_auth/PhUnits/array"
	"reflect"
	"strings"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"

	"encoding/json"
	"ph_auth/PhModel"
	"ph_auth/PhServer"
	"time"
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

	//TODO:存在Redis中的token信息[Scope信息]无具体操作权限,如果能Set进去,就不用查询数据库,毕竟页面验证比较频繁 => 不用访问数据库查询Scope?
	accRes := PhModel.Account{}
	accOut := PhModel.Account{}
	cond := bson.M{"_id": bson.ObjectIdHex(token.GetUserID())}
	err = h.db.FindOneByCondition(&accRes, &accOut, cond)
	if err != nil {
		panic(err.Error())
	}

	empModel := PhModel.Employee{}
	err = h.db.FindOneByCondition(&empModel, &empModel, bson.M{"_id": bson.ObjectIdHex(accOut.EmployeeID)})
	if err != nil {
		panic(err.Error())
	}

	scopeReq := api2go.Request{
		QueryParams: map[string][]string{
			"group-id": {empModel.GroupID},
		},
	}
	scopeIn := PhModel.Scope{}
	var scopeModels []PhModel.Scope
	err = h.db.FindMulti(scopeReq, &scopeIn, &scopeModels, -1, -1)
	if err == nil {
		for i, iter := range scopeModels {
			h.db.ResetIdWithId_(&iter)
			scopeModels[i] = iter
		}
	} else {
		panic(err.Error())
	}

	applyScope := token.GetScope()	// 申请Scope

	detailScope := strings.Split(applyScope, "/")
	//level := detailScope[0]       // Pharbers 官网 App 单个系统
	applyAccess := strings.Split(detailScope[1], ":")[0] // 申请的动作表述
	accessed := checkAccessScope(applyAccess, scopeModels)	//允许访问(判断是否有授权)
	indate := checkIndateScope(applyAccess, scopeModels)	//有效期内(判断授权是否过期)
	if accessed && indate {
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
	} else {
		panic(errors.ErrInvalidScope.Error())
	}
}

func (h PhTokenHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhTokenHandler) GetHandlerMethod() string {
	return h.Method
}

func singleAppGetScope(accScope []PhModel.Scope, applyScopes []string, prefix, level string) (bool, string) {
	var (
		scope     string
		scopeTemp map[string][]string
		temp      []string
	)
	scopeTemp = make(map[string][]string)
	for _, v := range accScope {
		temp = append(temp, v.Access+":"+v.Operation)
	}
	for _, applyScope := range applyScopes {
		if array.IsExistItem(prefix+":"+applyScope, temp) {
			scopeTemp[prefix] = append(scopeTemp[prefix], applyScope)
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

func checkAccessScope(applyAccess string, accScopes []PhModel.Scope) (accessed bool) {

	accessed = false
	for _, v := range accScopes {
		if v.Access == applyAccess {
			operations := strings.Split(v.Operation, "#")
			if checkOperationScope(operations[1]) {
				accessed = true
				return
			}
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

func checkOperationScope(operationCmd string) (allowed bool) {

	allowed = false
	operationCmdArr := strings.Split(operationCmd, "")
	if len(operationCmdArr) != 3 {
		panic("Scope OperationCmd Error!")
	}
	//TODO:针对不同情况验证权限[还需要再想想]
	if operationCmdArr[2] == "x" {
		allowed = true
	}
	return
}
