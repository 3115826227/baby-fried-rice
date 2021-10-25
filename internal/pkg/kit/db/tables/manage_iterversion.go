package tables

// 迭代版本信息
type IterativeVersion struct {
	Version string `gorm:"column:version;primaryKey"`
	Content string `gorm:"column:content;type:text"`
	// 发布状态 0-未发布，1-已发布
	Status          bool  `gorm:"column:status"`
	CreateTimestamp int64 `gorm:"column:create_timestamp"`
	UpdateTimestamp int64 `gorm:"column:update_timestamp"`
}

func (table *IterativeVersion) TableName() string {
	return "baby_manage_iter_version"
}
