package NtmFactory

import (
	"github.com/PharbersDeveloper/NtmPods/NtmDataStorage"
	"github.com/PharbersDeveloper/NtmPods/NtmHandler"
	"github.com/PharbersDeveloper/NtmPods/NtmModel"
	"github.com/PharbersDeveloper/NtmPods/NtmResource"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/PharbersDeveloper/NtmPods/NtmMiddleware"
)

type NtmTable struct{}

var NTM_MODEL_FACTORY = map[string]interface{}{
	"NtmImage":                NtmModel.Image{},
	"NtmPolicy":               NtmModel.Policy{},
	"NtmHospital":             NtmModel.Hospital{},
	"NtmDepartment":           NtmModel.Department{},
	"NtmRegion":               NtmModel.Region{},
	"NtmProduct":              NtmModel.Product{},
	"NtmProductConfig":        NtmModel.ProductConfig{},
	"NtmRepresentative":       NtmModel.Representative{},
	"NtmManagerConfig":        NtmModel.ManagerConfig{},
	"NtmRepresentativeConfig": NtmModel.RepresentativeConfig{},
	"NtmRegionConfig":         NtmModel.RegionConfig{},
	"NtmHospitalConfig":       NtmModel.HospitalConfig{},
	"NtmResourceConfig":       NtmModel.ResourceConfig{},
	"NtmGoodsConfig":          NtmModel.GoodsConfig{},
	"NtmBusinessInput":        NtmModel.BusinessInput{},
	"NtmRepresentativeInput":  NtmModel.RepresentativeInput{},
	"NtmManagerInput":         NtmModel.ManagerInput{},
	"NtmPaperInput":           NtmModel.PaperInput{},
	"NtmDestConfig":           NtmModel.DestConfig{},
	"NtmScenario":             NtmModel.Scenario{},
	"NtmProposal":             NtmModel.Proposal{},
	"NtmUseableProposal":      NtmModel.UseableProposal{},
	"NtmPaper":                NtmModel.Paper{},
}

var NTM_STORAGE_FACTORY = map[string]interface{}{
	"NtmImageStorage":                NtmDataStorage.NtmImageStorage{},
	"NtmPolicyStorage":               NtmDataStorage.NtmPolicyStorage{},
	"NtmHospitalStorage":             NtmDataStorage.NtmHospitalStorage{},
	"NtmDepartmentStorage":           NtmDataStorage.NtmDepartmentStorage{},
	"NtmRegionStorage":               NtmDataStorage.NtmRegionStorage{},
	"NtmProductStorage":              NtmDataStorage.NtmProductStorage{},
	"NtmProductConfigStorage":        NtmDataStorage.NtmProductConfigStorage{},
	"NtmRepresentativeStorage":       NtmDataStorage.NtmRepresentativeStorage{},
	"NtmManagerConfigStorage":        NtmDataStorage.NtmManagerConfigStorage{},
	"NtmRepresentativeConfigStorage": NtmDataStorage.NtmRepresentativeConfigStorage{},
	"NtmRegionConfigStorage":         NtmDataStorage.NtmRegionConfigStorage{},
	"NtmHospitalConfigStorage":       NtmDataStorage.NtmHospitalConfigStorage{},
	"NtmResourceConfigStorage":       NtmDataStorage.NtmResourceConfigStorage{},
	"NtmGoodsConfigStorage":          NtmDataStorage.NtmGoodsConfigStorage{},
	"NtmBusinessInputStorage":        NtmDataStorage.NtmBusinessInputStorage{},
	"NtmRepresentativeInputStorage":  NtmDataStorage.NtmRepresentativeInputStorage{},
	"NtmManagerInputStorage":         NtmDataStorage.NtmManagerInputStorage{},
	"NtmPaperInputStorage":           NtmDataStorage.NtmPaperInputStorage{},
	"NtmDestConfigStorage":           NtmDataStorage.NtmDestConfigStorage{},
	"NtmScenarioStorage":             NtmDataStorage.NtmScenarioStorage{},
	"NtmProposalStorage":             NtmDataStorage.NtmProposalStorage{},
	"NtmUseableProposalStorage":      NtmDataStorage.NtmUseableProposalStorage{},
	"NtmPaperStorage":                NtmDataStorage.NtmPaperStorage{},
}

var NTM_RESOURCE_FACTORY = map[string]interface{}{
	"NtmImageResource":                NtmResource.NtmImageResource{},
	"NtmPolicyResource":               NtmResource.NtmPolicyResource{},
	"NtmHospitalResource":             NtmResource.NtmHospitalResource{},
	"NtmDepartmentResource":           NtmResource.NtmDepartmentResource{},
	"NtmRegionResource":               NtmResource.NtmRegionResource{},
	"NtmProductResource":              NtmResource.NtmProductResource{},
	"NtmProductConfigResource":        NtmResource.NtmProductConfigResource{},
	"NtmRepresentativeResource":       NtmResource.NtmRepresentativeResource{},
	"NtmManagerConfigResource":        NtmResource.NtmManagerConfigResource{},
	"NtmRepresentativeConfigResource": NtmResource.NtmRepresentativeConfigResource{},
	"NtmRegionConfigResource":         NtmResource.NtmRegionConfigResource{},
	"NtmHospitalConfigResource":       NtmResource.NtmHospitalConfigResource{},
	"NtmResourceConfigResource":       NtmResource.NtmResourceConfigResource{},
	"NtmGoodsConfigResource":          NtmResource.NtmGoodsConfigResource{},
	"NtmBusinessInputResource":        NtmResource.NtmBusinessInputResource{},
	"NtmRepresentativeInputResource":  NtmResource.NtmRepresentativeInputResource{},
	"NtmManagerInputResource":         NtmResource.NtmManagerInputResource{},
	"NtmPaperInputResource":           NtmResource.NtmPaperInputResource{},
	"NtmDestConfigResource":           NtmResource.NtmDestConfigResource{},
	"NtmScenarioResource":             NtmResource.NtmScenarioResource{},
	"NtmProposalResource":             NtmResource.NtmProposalResource{},
	"NtmUseableProposalResource":      NtmResource.NtmUseableProposalResource{},
	"NtmPaperResource":                NtmResource.NtmPaperResource{},
}

var NTM_FUNCTION_FACTORY = map[string]interface{}{
	"NtmCommonPanicHandle":         	NtmHandler.CommonPanicHandle{},
	"NtmAccountHandler":            	NtmHandler.NtmAccountHandler{},
	"NtmGetUseableProposalsHandler": 	NtmHandler.NtmGetUseableProposalsHandler{},
	"NtmGeneratePaperHandler": 			NtmHandler.NtmGeneratePaperHandler{},
}
var NTM_MIDDLEWARE_FACTORY = map[string]interface{}{
	"NtmCheckTokenMiddleware": NtmMiddleware.NtmCheckTokenMiddleware{},
}

var NTM_DAEMON_FACTORY = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
}

func (t NtmTable) GetModelByName(name string) interface{} {
	return NTM_MODEL_FACTORY[name]
}

func (t NtmTable) GetResourceByName(name string) interface{} {
	return NTM_RESOURCE_FACTORY[name]
}

func (t NtmTable) GetStorageByName(name string) interface{} {
	return NTM_STORAGE_FACTORY[name]
}

func (t NtmTable) GetDaemonByName(name string) interface{} {
	return NTM_DAEMON_FACTORY[name]
}

func (t NtmTable) GetFunctionByName(name string) interface{} {
	return NTM_FUNCTION_FACTORY[name]
}

func (t NtmTable) GetMiddlewareByName(name string) interface{} {
	return NTM_MIDDLEWARE_FACTORY[name]
}
