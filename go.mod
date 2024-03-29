module baby-fried-rice

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)

require (
	github.com/Shopify/sarama v1.30.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibabacloud-go/darabonba-openapi v0.1.5
	github.com/alibabacloud-go/dysmsapi-20170525/v2 v2.0.1
	github.com/alibabacloud-go/tea v1.1.17
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/juju/testing v0.0.0-20211215003918-77eb13d6cad2 // indirect
	github.com/minio/minio-go/v7 v7.0.19
	github.com/nsqio/go-nsq v1.0.8
	github.com/olivere/elastic/v7 v7.0.30
	github.com/pingcap/tidb v0.0.0-20190108123336-c68ee7318319
	github.com/pion/rtcp v1.2.6
	github.com/pion/webrtc/v3 v3.0.32
	github.com/pkg/errors v0.9.1
	github.com/qiniu/api.v7/v7 v7.8.2
	github.com/satori/go.uuid v1.2.0
	github.com/siddontang/go-mysql-elasticsearch v0.0.0-20200822025838-fe261969558b
	github.com/swaggo/swag v1.7.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/ugorji/go v1.1.13 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.4 // indirect
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
)
