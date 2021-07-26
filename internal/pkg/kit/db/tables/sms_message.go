package tables

import "baby-fried-rice/internal/pkg/kit/constant"

// 短信发送记录
type SendMessageLog struct {
	ID            int64  `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	AccountId     string `gorm:"column:account_id" json:"account_id"`
	Phone         string `gorm:"column:phone" json:"phone"`
	Code          string `gorm:"column:code" json:"code"`
	SignName      string `gorm:"column:sign_name" json:"sign_name"`
	TemplateCode  string `gorm:"column:template_code" json:"template_code"`
	SendTimestamp int64  `gorm:"column:send_timestamp" json:"send_timestamp"`
}

func (table *SendMessageLog) TableName() string {
	return "baby_sms_send_message_log"
}

// 短信发送模板
type SendMessageTemplate struct {
	ID               int64                        `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	Name             string                       `gorm:"name" json:"name"`
	Code             constant.SmsTemplateCodeType `gorm:"column:code;unique" json:"code"`
	TemplateCode     string                       `gorm:"column:template_code;unique" json:"template_code"`
	Content          string                       `gorm:"column:content" json:"content"`
	CreatedTimestamp int64                        `gorm:"column:create_timestamp" json:"created_timestamp"`
}

func (table *SendMessageTemplate) TableName() string {
	return "baby_sms_send_message_template"
}
