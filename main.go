package main

import (
  "fmt"
  "os"
  "net/http"

  "github.com/alfredyang1986/BmServiceDef/BmConfig"
  "github.com/alfredyang1986/BmServiceDef/BmPodsDefine"
  "github.com/alfredyang1986/BmServiceDef/BmApiResolver"

  "github.com/PharbersDeveloper/PhAuthServer/PhFactory"
)

func main() {
  version := "v0"
  prodEnv := "NTM_HOME"
  fmt.Println("NTM pods archi begins, version =", version)

  fmt.Println("This is Pharbers Authorization Server")
  fac := NtmFactory.NtmTable{}
  var pod = BmPodsDefine.Pod{Name: "Pharbers Auth", Factory: fac}
  ntmHome := os.Getenv(prodEnv)
  pod.RegisterSerFromYAML(ntmHome + "/resource/service-def.yaml")

  var bmRouter BmConfig.BmRouterConfig
  bmRouter.GenerateConfig(prodEnv)

  addr := bmRouter.Host + ":" + bmRouter.Port
  fmt.Println("Listening on ", addr)
  api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
  pod.RegisterAllResource(api)

  pod.RegisterAllFunctions(version, api)
  pod.RegisterAllMiddleware(api)

  handler := api.Handler().(*httprouter.Router)
  pod.RegisterPanicHandler(handler)
  http.ListenAndServe(":"+bmRouter.Port, handler)

  fmt.Println("NTM pods archi ends, version =", version)
}
