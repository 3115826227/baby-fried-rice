package rsp

import "baby-fried-rice/internal/pkg/kit/db/tables"

func SmsLogModelToRsp(l tables.SendMessageLog) SmsLog {
	return SmsLog{
		Id:            l.ID,
		AccountId:     l.AccountId,
		Phone:         l.Phone,
		Code:          l.Code,
		SignName:      l.SignName,
		TemplateCode:  l.TemplateCode,
		SendTimestamp: l.SendTimestamp,
	}
}

type SmsLog struct {
	Id            int64  `json:"id"`
	AccountId     string `json:"account_id"`
	Phone         string `json:"phone"`
	Code          string `json:"code"`
	SignName      string `json:"sign_name"`
	TemplateCode  string `json:"template_code"`
	SendTimestamp int64  `json:"send_timestamp"`
}
