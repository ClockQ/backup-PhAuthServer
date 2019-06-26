package PhMiddleware

import (
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/manyminds/api2go"
	"net/http"
	"ph_auth/PhServer"
	"reflect"
	"strings"
)

var PhCheckToken PhCheckTokenMiddleware

type PhCheckTokenMiddleware struct {
	Args []string
	rd   *BmRedis.BmRedis
	md   *BmMongodb.BmMongodb
}

func (ctm PhCheckTokenMiddleware) NewCheckTokenMiddleware(args ...interface{}) PhCheckTokenMiddleware {
	var r *BmRedis.BmRedis
	var m *BmMongodb.BmMongodb
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
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	PhCheckToken = PhCheckTokenMiddleware{Args: ag, rd: r, md:m}
	return PhCheckToken
}

func (ctm PhCheckTokenMiddleware) DoMiddleware(c api2go.APIContexter, w http.ResponseWriter, r *http.Request) {
	// 垃圾的一批
	if !strings.Contains(r.RequestURI, "applyusers") {
		if err := ctm.CheckTokenFormFunction(w, r); err != nil {
			panic(err.Error())
		}
	}
}

func (ctm PhCheckTokenMiddleware) CheckTokenFormFunction(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Add("Content-Type", "application/json")

	sv := PhServer.GetInstance(ctm.md, ctm.rd)
	token, err := sv.ValidationBearerToken(r)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("PhCheckTokenMiddleware => token = ", token)
	return

}
