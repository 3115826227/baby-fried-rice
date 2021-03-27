package model

import "baby-fried-rice/internal/pkg/module/accountDao/model/tables"

type RespUserLogin struct {
	User   tables.AccountUser       `json:"user"`
	Detail tables.AccountUserDetail `json:"detail"`
}
