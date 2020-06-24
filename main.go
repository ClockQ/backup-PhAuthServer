package main

import (
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmApiResolver"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmServiceDef/BmPodsDefine"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"github.com/PharbersDeveloper/PhAuthServer/PhClient"
	"github.com/PharbersDeveloper/PhAuthServer/PhFactory"
)

func main() {

	version := os.Getenv("PH_OAUTH__VERSION")
	if version == "" {
		version = "v0"
	}

	serviceDef := os.Getenv("PH_OAUTH__SERVICE_DEF")
	if serviceDef == "" {
		serviceDef = "resources/service-def.yaml"
	}

	endPoint := os.Getenv("PH_OAUTH__END_POINT")
	if endPoint == "" {
		endPoint = "resources/endpoint.json"
	}

	host := os.Getenv("PH_OAUTH__HOST")
	port := os.Getenv("PH_OAUTH__PORT")
	if port == "" {
		port = "9096"
	}

	log.Println("Pharbers Auth Server begins, version =", version)

	fac := PhFactory.PhTable{}
	var pod = BmPodsDefine.Pod{Name: "Pharbers Auth", Factory: fac}
	pod.RegisterSerFromYAML(serviceDef)

	result := BmConfig.BmGetConfigMap(endPoint)
	fmt.Println(result)
	PhClient.EndPoint.RegisterEndPoint(result)


	addr := host + ":" + port
	log.Println("Pharbers Auth Server Listening on", addr)
	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)

	pod.RegisterAllFunctions(version, api)
	pod.RegisterAllMiddleware(api)

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	handler := api.Handler().(*httprouter.Router)

	pod.RegisterPanicHandler(handler)
	http.ListenAndServe(addr, c.Handler(handler))

	log.Println("Pharbers Auth Server begins, version =", version)
}
