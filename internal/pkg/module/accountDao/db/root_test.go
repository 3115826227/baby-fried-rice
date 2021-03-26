package db

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"testing"
	"time"
)

func TestAddRoot(t *testing.T) {
	now := time.Now()
	root := tables.AccountRoot{
		CommonField: tables.CommonField{
			ID:        handle.GenerateID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		LoginName: "root",
		//Username:  "超级管理员",
		Password:  handle.EncodePassword("root"),
	}
	if err := AddRoot(root); err != nil {
		panic(err)
	}
}
