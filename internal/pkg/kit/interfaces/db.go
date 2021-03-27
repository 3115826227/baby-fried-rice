package interfaces

import "gorm.io/gorm"

type DB interface {
	// 初始化建表操作
	InitTables(dos ...DataObject) error
	// 添加数据
	CreateObject(do DataObject) error
	// 判断数据是否存在
	ExistObject(query map[string]interface{}, do DataObject) (bool, error)
	// 删除数据
	DeleteObject(do DataObject) error
	// 更新数据
	UpdateObject(do DataObject) error
	// 查询单个数据
	GetObject(query map[string]interface{}, do DataObject) error
	// 获取db
	GetDB() *gorm.DB
	// 添加多条数据
	CreateMulti(bean ...interface{}) error
}
