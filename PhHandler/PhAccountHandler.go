package PhHandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"ph_auth/PhModel"
	"ph_auth/PhUnits/uuid"
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

func (h PhAccountHandler) ForgetPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	var parameter map[string]interface{}

	json.Unmarshal(body,&parameter)

	email, eOk := parameter["email"]

	if !eOk {
		response["status"] = "error"
		response["msg"] = "Email参数缺失"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	account := PhModel.Account{}
	var out PhModel.Account
	condition := bson.M{"email": email.(string)}

	err = h.db.FindOneByCondition(&account, &out, condition)
	if err != nil && err.Error() == "not found"{
		response["status"] = "error"
		response["msg"] = "未查询到该Email"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	uuid, _ := uuid.NewRandom()

	client := h.rd.GetRedisClient()
	defer client.Close()
	pipe := client.Pipeline()
	pipe.Set(uuid.String(), out.Email, time.Minute * 5)
	pipe.Exec()

	url := fmt.Sprint(h.Args[3],"/reset-password", "?uuid=", uuid.String(), "&email=", out.Email, "&progress=2")
	content := []byte(`{
		"email": "`+ out.Email +`",
		"subject": "申请修改密码",
		"content": "<a href=`+ url +`>点击修改密码</a>",
		"content-type": "text/html; charset=UTF-8"}`)

	h.sendEmail(r, content)


	response["status"] = "success"
	response["msg"] = "Email已发送"
	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 0
}

func (h PhAccountHandler) VerifyUUID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	var parameter map[string]interface{}

	json.Unmarshal(body,&parameter)

	uuid, uOk := parameter["uuid"]

	if !uOk {
		response["status"] = "error"
		response["msg"] = "uuid参数缺失"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	client := h.rd.GetRedisClient()
	defer client.Close()

	_, err = client.Get(uuid.(string)).Result()

	if err == redis.Nil {
		response["status"] = "error"
		response["msg"] = "超时"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	} else if err != nil {
		response["status"] = "error"
		response["msg"] = "未知错误"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}
	response["status"] = "success"
	response["msg"] = "验证成功"

	client.Del(uuid.(string))

	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 0
}

func (h PhAccountHandler) UpdatePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	var parameter map[string]interface{}

	json.Unmarshal(body,&parameter)

	password, pOk := parameter["password"]
	email, eOk := parameter["email"]

	if !pOk {
		response["status"] = "error"
		response["msg"] = "参数缺失"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	} else if !eOk {
		response["status"] = "error"
		response["msg"] = "Email参数缺失"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	account := PhModel.Account{}
	var out PhModel.Account

	condition := bson.M{"email": email.(string)}

	err = h.db.FindOneByCondition(&account, &out, condition)
	if err != nil && err.Error() == "not found"{
		response["status"] = "error"
		response["msg"] = "未查询到该Email"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	out.Password = password.(string)

	err = h.db.Update(&out)
	if err != nil {
		response["status"] = "error"
		response["msg"] = "更新失败"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	response["status"] = "success"
	response["msg"] = "更新成功"
	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 1
}

func (h PhAccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	// 注释掉token验证
	//token := PhMiddleware.PhCheckTokenMiddleware{Args: h.Args, Md: h.db, Rd: h.rd}
	//err := token.CheckTokenFormFunction(w, r)
	//if err != nil {
	//	panic(err.Error())
	//}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var parameter map[string]string
	err = json.Unmarshal(body, &parameter)
	if err != nil {
		panic(err.Error())
	}

	accountIn := PhModel.Account{ID: parameter["account"]}
	accountOut := PhModel.Account{}
	h.db.FindOne(&accountIn, &accountOut)

	if accountOut.ID == "" {
		panic("Account is Error")
	}

	employeeIn := PhModel.Employee{ID: accountOut.EmployeeID}
	employeeOut := PhModel.Employee{ID: accountOut.EmployeeID}
	h.db.FindOne(&employeeIn, &employeeOut)

	var employeeOuts []PhModel.Employee
	p := api2go.Request{
		QueryParams: map[string][]string{
			"group-ids": {employeeOut.GroupID},
		},
	}
	err = h.db.FindMulti(p, &PhModel.Employee{}, &employeeOuts, -1, -1)
	if err == nil {
		for i, iter := range employeeOuts {
			h.db.ResetIdWithId_(&iter)
			employeeOuts[i] = iter
		}
	} else {
		panic(err.Error())
	}

	var employeeIds []string
	for _, v := range employeeOuts {
		employeeIds = append(employeeIds, v.ID)
	}

	p = api2go.Request{
		QueryParams: map[string][]string{
			"employee-ids": employeeIds,
		},
	}
	var accountOuts []PhModel.Account
	err = h.db.FindMulti(p, &PhModel.Account{}, &accountOuts, -1, -1)
	if err == nil {
		for i, iter := range accountOuts {
			h.db.ResetIdWithId_(&iter)
			accountOuts[i] = iter
		}
	} else {
		panic(err.Error())
	}

	var accountIds []string
	for _, v := range accountOuts {
		accountIds = append(accountIds, v.ID)
	}

	enc := json.NewEncoder(w)
	enc.Encode(accountIds)
	return 0
}

func (h PhAccountHandler) GetAccountNameById(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var parameter map[string]string
	err = json.Unmarshal(body, &parameter)
	if err != nil {
		panic(err.Error())
	}

	accountIn := PhModel.Account{ID: parameter["account"]}
	accountOut := PhModel.Account{}
	h.db.FindOne(&accountIn, &accountOut)

	if accountOut.ID == "" {
		panic("Account is Error")
	}

	result := map[string]string{
		"accountId": 	accountOut.ID,
		"accountName": 	accountOut.Username,
	}

	enc := json.NewEncoder(w)
	enc.Encode(result)

	return 0
}

func (h PhAccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAccountHandler) GetHandlerMethod() string {
	return h.Method
}

func (h PhAccountHandler) sendEmail(r *http.Request, content []byte) {
	// 拼接转发的URL
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	resource := fmt.Sprint(h.Args[1], "/", h.Args[0], "/", h.Args[2])
	mergeURL := strings.Join([]string{scheme, resource}, "")

	fmt.Println(mergeURL)

	// 转发
	client := &http.Client{}
	req, _ := http.NewRequest("POST", mergeURL, bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/json")

	for k, v := range r.Header {
		req.Header.Add(k, v[0])
	}

	client.Do(req)

}
