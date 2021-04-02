package tables

import "time"

type CommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updated_at"`
}

type CommonIntField struct {
	ID        int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updated_at"`
}
