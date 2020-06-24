package PhFactory

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/PharbersDeveloper/PhAuthServer/PhMiddleware"

	"github.com/PharbersDeveloper/PhAuthServer/PhDataStorage"
	"github.com/PharbersDeveloper/PhAuthServer/PhHandler"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"github.com/PharbersDeveloper/PhAuthServer/PhResource"
)

type PhTable struct{}

var PhModelFactory = map[string]interface{}{
	"PhAccount": PhModel.Account{},
	"PhEmployee": PhModel.Employee{},
	"PhGroup": PhModel.Group{},
	"PhCompany": PhModel.Company{},
	"PhRole": PhModel.Role{},
	"PhApplyuser": PhModel.Applyuser{},
}

var PhStorageFactory = map[string]interface{}{
	"PhAccountStorage": PhDataStorage.PhAccountStorage{},
	"PhEmployeeStorage": PhDataStorage.PhEmployeeStorage{},
	"PhGroupStorage": PhDataStorage.PhGroupStorage{},
	"PhCompanyStorage": PhDataStorage.PhCompanyStorage{},
	"PhRoleStorage": PhDataStorage.PhRoleStorage{},
	"PhApplyuserStorage": PhDataStorage.PhApplyuserStorage{},
}

var PhResourceFactory = map[string]interface{}{
	"PhAccountResource": PhResource.PhAccountResource{},
	"PhEmployeeResource": PhResource.PhEmployeeResource{},
	"PhGroupResource": PhResource.PhGroupResource{},
	"PhCompanyResource": PhResource.PhCompanyResource{},
	"PhRoleResource": PhResource.PhRoleResource{},
	"PhApplyuserResource": PhResource.PhApplyuserResource{},
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
	"PhForgetPasswordHandler": 		PhHandler.PhAccountHandler{},
	"PhVerifyUUIDHandler": 			PhHandler.PhAccountHandler{},
	"PhUpdatePasswordHandler": 		PhHandler.PhAccountHandler{},
	"PhGetAccountsHandler": 		PhHandler.PhAccountHandler{},
	"PhGetAccountNameByIdHandler": 		PhHandler.PhAccountHandler{},

}
var PhMiddlewareFactory = map[string]interface{}{
	"PhCheckTokenMiddleware": PhMiddleware.PhCheckTokenMiddleware{},
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
