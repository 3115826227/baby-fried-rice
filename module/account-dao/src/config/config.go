package config

import (
	"github.com/jinzhu/configor"
	"os"
)

var Config = struct {
	MysqlUrl string `env:"MYSQL_URL" required:"true"`
}{}

var Permission struct {
	Permission interface{} `env:"PERMISSION" required:"true"`
}

const (
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"

	DefaultRoleName      = "默认"
	AdminRoleName        = "管理员"
	RootSchoolOrganizeId = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa"

	StudentVerifyTrue  = "已认证"
	StudentVerifyFalse = "未认证"

	DefaultSubAdminPassword = "1234"
	DefaultAdminPassword    = "1234"
	DefaultUserEncryMd5     = "md5"

	GetIpAddress = "http://ip-api.com/json"
)

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/account-dao"

func init() {
	var err error
	if err = configor.Load(&Config, "etc/config.yaml"); err != nil {
		panic(err)
	}
	if err = configor.Load(&Permission, "etc/permission.yaml"); err != nil {
		panic(err)
	}
}
