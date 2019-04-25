module ph_auth

replace golang.org/x/text => github.com/golang/text v0.3.0

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190422165155-953cdadca894

replace golang.org/x/net => github.com/golang/net v0.0.0-20190424024250-574d568418ea

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190422183909-d864b10871cd

replace golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190402181905-9f3314589c9a

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190328211700-ab21143f2384

replace cloud.google.com/go => github.com/googleapis/google-cloud-go v0.34.0

replace google.golang.org/appengine => github.com/golang/appengine v1.4.0

go 1.12

require (
	github.com/PharbersDeveloper/PhAuthServer v0.0.0-20190424041951-c99c688efad6
	github.com/alfredyang1986/BmServiceDef v0.0.0-20190410064517-b341f9e1e85e
	github.com/alfredyang1986/blackmirror v0.0.0-20190305121812-d8d7643fb552 // indirect
	github.com/gavv/httpexpect v1.0.0 // indirect
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813 // indirect
	github.com/go-redis/redis v6.15.2+incompatible // indirect
	github.com/julienschmidt/httprouter v1.2.0
	github.com/manyminds/api2go v0.0.0-20190324173508-d4f7fae65b4b
	github.com/rs/cors v1.6.0
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	golang.org/x/oauth2 v0.0.0-00010101000000-000000000000
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/oauth2.v3 v3.10.0
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
