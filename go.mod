module baby-fried-rice

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)

require (
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/onsi/ginkgo v1.15.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/qiniu/api.v7/v7 v7.8.2
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
)
