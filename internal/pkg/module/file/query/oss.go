package query

import (
	"baby-fried-rice/internal/pkg/module/file/db"
	"baby-fried-rice/internal/pkg/module/file/model/tables"
)

func GetOssMeta(id int) (meta tables.OssMeta, err error) {
	err = db.GetDB().GetObject(map[string]interface{}{"id": id}, &meta)
	return
}
