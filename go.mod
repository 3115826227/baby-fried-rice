module baby-fried-rice

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibabacloud-go/darabonba-openapi v0.1.5
	github.com/alibabacloud-go/dysmsapi-20170525/v2 v2.0.1
	github.com/alibabacloud-go/tea v1.1.17
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.0
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/nsqio/go-nsq v1.0.8
	github.com/onsi/ginkgo v1.15.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/qiniu/api.v7/v7 v7.8.2
	github.com/satori/go.uuid v1.2.0
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.4 // indirect
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
)
