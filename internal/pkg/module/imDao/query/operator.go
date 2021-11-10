package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/module/imDao/db"
	"gorm.io/gorm"
)

func GetOperator() {

}

func ExistOperator(origin, receive string, optType im.OptType) (bool, error) {
	var opt tables.Operator
	if err := db.GetDB().GetObject(map[string]interface{}{
		"origin":   origin,
		"receive":  receive,
		"opt_type": optType,
	}, &opt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if opt.Confirm == 2 {
		// 如果已经拒绝，则可以当做不存在
		return false, nil
	}
	return true, nil
}
