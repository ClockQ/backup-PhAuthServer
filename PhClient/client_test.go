package PhClient

import (
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"os"
	"testing"
)

func TestClientEndPoint(t *testing.T) {
	prodEnv := os.Getenv("PH_AUTH_HOME")
	fmt.Println(prodEnv)
	result := BmConfig.BmGetConfigMap(prodEnv + "/resource/" + "endpoint.json")
	fmt.Println(result)
	EndPoint.RegisterEndPoint(result)
	EndPoint.ConfigFromURIParameter(nil)
}
