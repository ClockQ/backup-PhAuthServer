package PhHandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"os"
)

type PhLoginPageHandler struct {
	Method     string
	HttpMethod string
	Args       []string
}

func (h PhLoginPageHandler) NewLoginPageHandler(args ...interface{}) PhLoginPageHandler {
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
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

	return PhLoginPageHandler{Method: md, HttpMethod: hm, Args: ag}
}

func (h PhLoginPageHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	file, err := os.Open(h.Args[0])
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return 1
	}
	fi, _ := file.Stat()
	http.ServeContent(w, r, file.Name(), fi.ModTime(), file)
	return 0
}

func (h PhLoginPageHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhLoginPageHandler) GetHandlerMethod() string {
	return h.Method
}
