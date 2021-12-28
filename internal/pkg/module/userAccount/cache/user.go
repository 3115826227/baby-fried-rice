package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func NewUserPhoneCodeKey(accountId string) string {
	return fmt.Sprintf("%v:%v", constant.AccountPhoneVerifyCodePrefix, accountId)
}

func SetUserPhoneCode(accountId, code string) error {
	return GetCache().GetRedis().Set(NewUserPhoneCodeKey(accountId), code, time.Minute).Err()
}

func GetUserPhoneCode(accountId string) (string, bool, error) {
	code, err := GetCache().Get(NewUserPhoneCodeKey(accountId))
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return code, true, nil
}
