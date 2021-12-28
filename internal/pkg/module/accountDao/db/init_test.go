package db

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"strconv"
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {
	var conf = models.Mysql{
		Host:     "127.0.0.1",
		Port:     23306,
		Username: "root",
		Password: "123456",
		DBName:   "baby-account",
	}
	if err := InitDB(conf.GetMysqlUrl()); err != nil {
		panic(err)
	}

	provider := &db.BatchProvider{
		TableName:      (&tables.AccountUser{}).TableName(),
		Fields:         []string{"id", "account_id", "login_name", "password", "encode_type", "create_time", "update_time"},
		ConflictFields: []string{"login_name"},
		UpdateFields:   []string{"password", "update_time"},
		BatchAmount:    1000,
	}

	var number = 50000
	var record = make([][]interface{}, 0)
	for i := 1; i <= number; i++ {
		var now = time.Now()
		var accountId = handle.GenerateSerialNumber()
		record = append(record, []interface{}{
			handle.GenerateID(), accountId,
			"test-" + strconv.Itoa(i), handle.EncodePassword(accountId, "123456"),
			constant.DefaultUserEncryMd5, now, now,
		})
	}
	if err := provider.Update(GetDB().GetDB(), record); err != nil {
		panic(err)
	}

	detailProvider := &db.BatchProvider{
		TableName:      (&tables.AccountUserDetail{}).TableName(),
		Fields:         []string{"id", "account_id", "username", "create_time", "update_time"},
		ConflictFields: []string{"account_id"},
		UpdateFields:   []string{"username", "update_time"},
		BatchAmount:    1000,
	}
	var detailRecords = make([][]interface{}, 0)
	for i := 1; i <= number; i++ {
		var now = time.Now()
		detailRecords = append(detailRecords, []interface{}{
			record[i-1][0], record[i-1][1],
			"测试用户",
			now, now,
		})
	}
	if err := detailProvider.Update(GetDB().GetDB(), detailRecords); err != nil {
		panic(err)
	}

}
