package db

import (
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"go.uber.org/zap"
	"strings"
)

// 批量 insert 到 postgres
// updateFields 不为空时，conflict 时可触发更新；否则不更新
func load(tableName string, fields, conflictFields, updateFields []string, records [][]interface{}) (int, error) {
	cursor := DB
	valueNames := strings.Join(fields, ", ")
	args := make([]interface{}, 1)

	var valuePlaceHolder = strings.Repeat("?,", len(fields))
	valuePlaceHolder = "(" + valuePlaceHolder[:len(valuePlaceHolder)-1] + "),"
	valuePlaceHolders := strings.Repeat(valuePlaceHolder, len(records))
	valuePlaceHolders = valuePlaceHolders[:len(valuePlaceHolders)-1]
	for _, record := range records {
		args = append(args, record...)
	}

	sql := "insert into " + tableName + " (" + valueNames + ") values" + valuePlaceHolders
	if len(conflictFields) > 0 {
		onDups := make([]string, 0)
		sql += " on conflict(" + strings.Join(conflictFields, ", ") + ") do "
		if len(updateFields) > 0 {
			for _, field := range updateFields {
				onDups = append(onDups, field+"=excluded."+field)
			}
			sql += " update set " + strings.Join(onDups, ", ")
		} else {
			sql += " nothing"
		}
	}

	args[0] = sql
	rowsAffected := cursor.Exec(sql, args[1:]...).RowsAffected
	if rowsAffected != int64(len(records)) {
		log.Logger.Info("Postgres", zap.String("table name", tableName), zap.Int("to insert", len(records)),
			zap.Int64("inserted", rowsAffected))
	}
	return len(records), nil
}

// 按 BatchLoadAmount 的值批量导入 postgres
func Load(tableName string, fields, conflictFields, updateFields []string, records [][]interface{}) (int, error) {
	bulkSize := config.BatchLoadAmount
	index := 0
	success := 0
	var count int
	var end int
	var err error
	for index < len(records) {
		end = index + int(bulkSize)
		if end > len(records) {
			end = len(records)
		}
		count, err = load(tableName, fields, conflictFields, updateFields, records[index:end])
		success += count
		index = end
	}
	return success, err
}
