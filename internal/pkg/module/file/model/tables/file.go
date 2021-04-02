package tables

import "baby-fried-rice/internal/pkg/kit/models/tables"

// 文件信息
type File struct {
	tables.CommonField
	// 文件创建者
	Origin string `gorm:"column:origin" json:"origin"`
	// 文件权限类型： 1-public 2-protect 3-private
	PermissionType int `gorm:"column:permission_type;" json:"permission_type"`
	// 文件名
	FileName string `gorm:"column:file_name" json:"file_name"`
	// 文件类型：1-图片 2-文本 3-视频
	FileType int `gorm:"column:file_type" json:"file_type"`
	// 文件大小
	FileSize int64 `gorm:"column:file_size" json:"file_size"`
	// 文件存储地址
	Path string `gorm:"column:path" json:"path"`
}

func (table *File) TableName() string {
	return "baby_file"
}

func (table *File) Get() interface{} {
	return *table
}

// 文件组信息
type FileGroup struct {
	tables.CommonField
	// 文件组名
	Name string `gorm:"column:name" json:"name"`
	// 文件组描述信息
	Desc string `gorm:"column:desc" json:"desc"`
}

func (table *FileGroup) TableName() string {
	return "baby_file_group"
}

func (table *FileGroup) Get() interface{} {
	return *table
}

// 文件组关联信息
type FileGroupUserRelation struct {
	// 文件组id
	GroupID string `gorm:"column:group_id;unique_index:idx_group_user_id"`
	// 用户id
	UserID string `gorm:"column:user_id;unique_index:idx_group_user_id"`
}

func (table *FileGroupUserRelation) TableName() string {
	return "baby_file_group_user_rel"
}

func (table *FileGroupUserRelation) Get() interface{} {
	return *table
}
