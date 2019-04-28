package PhFactory

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"

	"ph_auth/PhDataStorage"
	"ph_auth/PhHandler"
	"ph_auth/PhModel"
	"ph_auth/PhResource"
)

type PhTable struct{}

var PhModelFactory = map[string]interface{}{
	"PhAccount": PhModel.Account{},
}

var PhStorageFactory = map[string]interface{}{
	"PhAccountStorage": PhDataStorage.PhAccountStorage{},
}

var PhResourceFactory = map[string]interface{}{
	"PhAccountResource": PhResource.PhAccountResource{},
}

var PhFunctionFactory = map[string]interface{}{
	"PhCommonPanicHandle":          PhHandler.CommonPanicHandle{},
	"PhLoginPageHandler":           PhHandler.PhLoginPageHandler{},
	"PhAuthPageHandler":            PhHandler.PhAuthPageHandler{},
	"PhAccountHandler":             PhHandler.PhAccountHandler{},
	"PhTokenHandler":               PhHandler.PhTokenHandler{},
	"PhTokenValidationHandler":     PhHandler.PhTokenHandler{},
	"PhAuthorizeHandler":           PhHandler.PhAuthorizeHandler{},
	"PhUserAgentHandler":           PhHandler.PhUserAgentHandler{},
	"PhGenerateAccessTokenHandler": PhHandler.PhAuthorizeHandler{},
	"PhRefreshAccessTokenHandler":  PhHandler.PhAuthorizeHandler{},
	"PhPasswordLoginHandler":       PhHandler.PhAuthorizeHandler{},
}
var PhMiddlewareFactory = map[string]interface{}{
	//"NtmCheckTokenMiddleware": NtmMiddleware.NtmCheckTokenMiddleware{},
}

var PhDaemonFactory = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
}

func (t PhTable) GetModelByName(name string) interface{} {
	return PhModelFactory[name]
}

func (t PhTable) GetResourceByName(name string) interface{} {
	return PhResourceFactory[name]
}

func (t PhTable) GetStorageByName(name string) interface{} {
	return PhStorageFactory[name]
}

func (t PhTable) GetDaemonByName(name string) interface{} {
	return PhDaemonFactory[name]
}

func (t PhTable) GetFunctionByName(name string) interface{} {
	return PhFunctionFactory[name]
}

func (t PhTable) GetMiddlewareByName(name string) interface{} {
	return PhMiddlewareFactory[name]
}
