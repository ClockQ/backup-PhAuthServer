package PhHandler

import (
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/url"
	"ph_auth/PhClient"
	"reflect"
	"strings"
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
	queryForm, _ := url.ParseQuery(r.URL.RawQuery)

	config := PhClient.EndPoint.ConfigFromURIParameter(r)
	redirectUrl := config.AuthCodeURL("xyz")
	if v := queryForm["status"]; len(v) > 0 && v[0] == "self" {
		redirectUrl += "&status=" + v[0]

		// 转发
		client := &http.Client{}
		req, _ := http.NewRequest("GET", redirectUrl, nil)
		for k, v := range r.Header {
			req.Header.Add(k, v[0])
		}
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Login Error")
		}
		data, err := ioutil.ReadAll(response.Body)

		str := string(data)
		preStr := str[:strings.LastIndex(str, "</body>")]
		sufStr := str[strings.LastIndex(str, "</body>"):]
		truncationUrl := redirectUrl[strings.LastIndex(redirectUrl, "?"):]
		insertContent := fmt.Sprint("<input type='hidden' id='parameter'", "value='", truncationUrl, "'", "/>")

		content := []byte(fmt.Sprint(preStr, insertContent, sufStr))

		w.Write(content)

		return 0
	}
	http.Redirect(w, r, redirectUrl, http.StatusFound)
	return 0
}

func (h PhUserAgentHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhUserAgentHandler) GetHandlerMethod() string {
	return h.Method
}
