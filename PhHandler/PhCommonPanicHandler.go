package PhHandler

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type CommonPanicHandle struct {
}

func (ctm CommonPanicHandle) NewCommonPanicHandle(args ...interface{}) CommonPanicHandle {
	return CommonPanicHandle{}
}

func (ctm CommonPanicHandle) HandlePanic(rw http.ResponseWriter, r *http.Request, p interface{}) {
	fmt.Println("CommonHandlePanic接收到", p)

	status := http.StatusOK

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")
	rw.WriteHeader(status)

	data := make(map[string]interface{})
	data["error"] = "invalid_grant"
	data["error_description"] = p
	json.NewEncoder(rw).Encode(data)

	//PhPanic.ErrInstance().ErrorReval(p.(string), rw)
}
