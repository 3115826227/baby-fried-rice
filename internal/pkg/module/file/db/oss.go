package db

import (
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/file/model/tables"
)

func AddOssMeta(ossMeta tables.OssMeta) (tables.OssMeta, error) {
	err := db.GetDB().CreateObject(&ossMeta)
	return ossMeta, err
}
