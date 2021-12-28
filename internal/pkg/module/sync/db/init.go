package db

//var (
//	accountClient interfaces.DB
//	shopClient    interfaces.DB
//	smsClient     interfaces.DB
//	spaceClient   interfaces.DB
//	imClient      interfaces.DB
//)
//
//func GetAccountDB() interfaces.DB {
//	return accountClient
//}
//
//func GetShopDB() interfaces.DB {
//	return shopClient
//}
//
//func GetSmsDB() interfaces.DB {
//	return smsClient
//}
//
//func GetSpaceDB() interfaces.DB {
//	return spaceClient
//}
//
//func GetImDB() interfaces.DB {
//	return imClient
//}
//
//func InitDB() (err error) {
//	var conf = config.GetConfig()
//	accountClient, err = db.NewClientDB(conf.Database.SubDatabase.AccountDatabase.GetMysqlUrl(), log.Logger)
//	if err != nil {
//		return
//	}
//	shopClient, err = db.NewClientDB(conf.Database.SubDatabase.ShopDatabase.GetMysqlUrl(), log.Logger)
//	if err != nil {
//		return
//	}
//	smsClient, err = db.NewClientDB(conf.Database.SubDatabase.SmsDatabase.GetMysqlUrl(), log.Logger)
//	if err != nil {
//		return
//	}
//	spaceClient, err = db.NewClientDB(conf.Database.SubDatabase.SpaceDatabase.GetMysqlUrl(), log.Logger)
//	if err != nil {
//		return
//	}
//	imClient, err = db.NewClientDB(conf.Database.SubDatabase.ImDatabase.GetMysqlUrl(), log.Logger)
//	if err != nil {
//		return
//	}
//	return
//}
