package PhPanic

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-uuid"
	"github.com/manyminds/api2go"
	"net/http"
	"strconv"
	"sync"
)

var PH_ERROR = errors.New("Pharbers error")

type tPhrror struct {
	m map[string]api2go.HTTPError
}

var e *tPhrror
var o sync.Once

//TODO:error definition load from file
func ErrInstance() *tPhrror {
	o.Do(func() {
		e = &tPhrror{
			m: map[string]api2go.HTTPError{
				"Auth Failed!": api2go.HTTPError{Errors: []api2go.Error{
					{
						Links: &api2go.ErrorLinks{
							About: "http://login",
						},
						Status: "401",
						Code:   "001",
						Title:  "Auth error!",
						Detail: "Auth error!",
						Source: &api2go.ErrorSource{
							Pointer: "#titleField",
						},
						Meta: map[string]interface{}{
							"creator": "jeorch",
						},
					},
				}},
				"no defind error!": api2go.HTTPError{Errors: []api2go.Error{
					{
						Links: &api2go.ErrorLinks{
							About: "http://404",
						},
						Status: "404",
						Code:   "9999",
						Title:  "no defind error!",
						Detail: "no defind error!",
						Source: &api2go.ErrorSource{
							Pointer: "#titleField",
						},
						Meta: map[string]interface{}{
							"creator": "jeorch",
						},
					},
				}},
			},
		}
	})
	return e
}

func (e *tPhrror) IsErrorDefined(ec string) bool {
	for k := range e.m {
		if k == ec {
			return true
		}
	}
	return false
}

func resetlHTTPErrorID(input api2go.HTTPError) {
	if len(input.Errors) != 0 {
		for i := range input.Errors {
			eid, _ := uuid.GenerateUUID()
			input.Errors[i].ID = eid
		}
	}
}

func (e *tPhrror) ErrorReval(err interface{}, w http.ResponseWriter) {
	es := err.(string)
	var hr api2go.HTTPError
	if e.IsErrorDefined(es) {
		hr = e.m[es]
	} else {
		hr = e.m["no defind error!"]
		hr.Errors[0].Detail = es
	}
	resetlHTTPErrorID(hr)
	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	statusCode, _ := strconv.Atoi(hr.Errors[0].Status)
	w.WriteHeader(statusCode)
	enc.Encode(hr)
}
