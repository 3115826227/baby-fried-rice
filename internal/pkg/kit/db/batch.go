package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type DialectorType string

const (
	DIALECTOR_MYSQL DialectorType = "mysql"
)

type BatchProvider struct {
	TableName      string   `json:"table_name"`
	Fields         []string `json:"fields"`
	ConflictFields []string `json:"conflict_fields, omitempty"`
	UpdateFields   []string `json:"update_fields"`
	BatchAmount    int      `json:"batch_amount"`
}

func (provider *BatchProvider) Update(engine *gorm.DB, records [][]interface{}) error {
	var (
		index = 0
		end   int
		err   error
	)
	for index < len(records) {
		end = index + provider.BatchAmount
		if end > len(records) {
			end = len(records)
		}
		if err = provider.load(engine, records[index:end]); err != nil {
			return err
		}
		index = end
	}
	return err
}

func (provider *BatchProvider) engineJudge(engine gorm.DB) DialectorType {
	switch engine.Dialector.Name() {
	case (&mysql.Dialector{}).Name():
		return DIALECTOR_MYSQL
	default:
		return ""
	}
}

func (provider *BatchProvider) constructSQL(records [][]interface{}, dialectorType DialectorType) (string, error) {
	switch dialectorType {
	case DIALECTOR_MYSQL:
		return provider.constructMYSQL(records), nil
	default:
		return "", fmt.Errorf("dialector type is invalid")
	}
}

func (provider *BatchProvider) constructMYSQL(records [][]interface{}) string {
	var (
		valueNames        string
		valuePlaceHolder  string
		valuePlaceHolders string
		sql               string
	)
	valueNames = strings.Join(provider.Fields, ", ")
	valuePlaceHolder = strings.Repeat("?,", len(provider.Fields))
	valuePlaceHolder = "(" + valuePlaceHolder[:len(valuePlaceHolder)-1] + "),"
	valuePlaceHolders = strings.Repeat(valuePlaceHolder, len(records))
	valuePlaceHolders = valuePlaceHolders[:len(valuePlaceHolders)-1]
	sql = "insert into " + provider.TableName + " (" + valueNames + ") values" + valuePlaceHolders
	var onDups []string
	sql += " on duplicate key "
	if len(provider.UpdateFields) > 0 {
		for _, field := range provider.UpdateFields {
			onDups = append(onDups, field+"=values("+field+")")
		}
		sql += "update " + strings.Join(onDups, ", ")
	} else {
		sql += "nothing"
	}
	return sql
}

func (provider *BatchProvider) load(engine *gorm.DB, records [][]interface{}) error {
	// 定义变量
	var (
		sql  string
		args []interface{}
		err  error
	)
	// 构造sql
	sql, err = provider.constructSQL(records, provider.engineJudge(*engine))
	if err != nil {
		return err
	}
	// 添加值列表
	for _, record := range records {
		args = append(args, record...)
	}
	return engine.Exec(sql, args...).Error
}
