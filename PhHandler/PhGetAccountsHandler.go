package PhHandler

import (
	"encoding/json"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"io/ioutil"
	"net/http"
	"ph_auth/PhModel"
	"reflect"
)

type PhGetAccountsHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h PhGetAccountsHandler) NewGetAccountsHandle(args ...interface{}) PhGetAccountsHandler {
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

	return PhGetAccountsHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h PhGetAccountsHandler) GetAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
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

	employeeIn := PhModel.Employee{ID: accountOut.EmployeeID}
	employeeOut := PhModel.Employee{ID: accountOut.EmployeeID}
	h.db.FindOne(&employeeIn, &employeeOut)

	var employeeOuts []PhModel.Employee
	p := api2go.Request{
		QueryParams: map[string][]string {
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
	for  _, v := range employeeOuts {
		employeeIds = append(employeeIds, v.ID)
	}

	p = api2go.Request{
		QueryParams: map[string][]string {
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
	for  _, v := range accountOuts {
		accountIds = append(accountIds, v.ID)
	}

	enc := json.NewEncoder(w)
	enc.Encode(accountIds)
	return 0
}

func (h PhGetAccountsHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhGetAccountsHandler) GetHandlerMethod() string {
	return h.Method
}
