package db

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/manage/config"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"time"
)

var (
	accountClient interfaces.DB
	shopClient    interfaces.DB
	smsClient     interfaces.DB
	spaceClient   interfaces.DB
	imClient      interfaces.DB
)

func initAdmin() {
	var adminId = "100000"
	var admin = tables.AccountAdmin{
		LoginName:  "root",
		Username:   "后台管理账号",
		Password:   handle.EncodePassword(adminId, "root1234"),
		EncodeType: "md5",
	}
	admin.ID = adminId
	now := time.Now()
	admin.CreatedAt, admin.UpdatedAt = now, now
	if exist, err := accountClient.ExistObject(map[string]interface{}{"id": admin.ID}, &tables.AccountAdmin{}); err != nil {
		return
	} else if !exist {
		if err = accountClient.CreateObject(&admin); err != nil {
			panic(err)
		}
	}
	return
}

func initSmsTemplate() {
	var registerTemplate = tables.SendMessageTemplate{
		ID:               1,
		Name:             "宝宝煎米果平台注册",
		Code:             constant.SmsRegisterCode,
		TemplateCode:     "SMS_164100177",
		Content:          "您正在申请手机注册，验证码为：${code}，5分钟内有效！",
		CreatedTimestamp: time.Now().Unix(),
	}
	if exist, err := smsClient.ExistObject(map[string]interface{}{"id": 1}, &tables.SendMessageTemplate{}); err != nil {
		panic(err)
	} else if !exist {
		if err = smsClient.CreateObject(&registerTemplate); err != nil {
			panic(err)
		}
	}
}

func GetAccountDB() interfaces.DB {
	return accountClient
}

func GetShopDB() interfaces.DB {
	return shopClient
}

func GetSmsDB() interfaces.DB {
	return smsClient
}

func GetSpaceDB() interfaces.DB {
	return spaceClient
}

func GetImDB() interfaces.DB {
	return imClient
}

func InitDB() (err error) {
	var conf = config.GetConfig()
	accountClient, err = db.NewClientDB(conf.Database.SubDatabase.AccountDatabase.GetMysqlUrl(), log.Logger)
	if err != nil {
		return
	}
	if err = accountClient.InitTables(
		&tables.AccountAdmin{},
		&tables.AccountAdminLoginLog{},
		&tables.IterativeVersion{},
		&tables.Communication{},
		&tables.CommunicationDetail{},
	); err != nil {
		panic(err)
	}
	initAdmin()
	shopClient, err = db.NewClientDB(conf.Database.SubDatabase.ShopDatabase.GetMysqlUrl(), log.Logger)
	if err != nil {
		return
	}
	if err = shopClient.InitTables(
		&tables.Commodity{},
		&tables.CommodityImageRel{},
	); err != nil {
		panic(err)
	}
	smsClient, err = db.NewClientDB(conf.Database.SubDatabase.SmsDatabase.GetMysqlUrl(), log.Logger)
	if err != nil {
		return
	}
	if err = smsClient.InitTables(
		&tables.SendMessageLog{},
		&tables.SendMessageTemplate{},
	); err != nil {
		panic(err)
	}
	spaceClient, err = db.NewClientDB(conf.Database.SubDatabase.SpaceDatabase.GetMysqlUrl(), log.Logger)
	if err != nil {
		return
	}
	imClient, err = db.NewClientDB(conf.Database.SubDatabase.ImDatabase.GetMysqlUrl(), log.Logger)
	if err != nil {
		return
	}
	initSmsTemplate()
	return
}
