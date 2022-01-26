package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func NewUserPhoneCodeKey(accountId string) string {
	return fmt.Sprintf("%v:%v", constant.AccountPhoneVerifyCodePrefix, accountId)
}

func SetUserPhoneCode(phoneCode models.UserPhoneCode) error {
	data, err := json.Marshal(phoneCode)
	if err != nil {
		return err
	}
	return GetCache().GetRedis().Set(NewUserPhoneCodeKey(phoneCode.AccountId), string(data), time.Minute).Err()
}

func GetUserPhoneCode(accountId string) (models.UserPhoneCode, bool, error) {
	data, err := GetCache().Get(NewUserPhoneCodeKey(accountId))
	if err != nil {
		if err == redis.Nil {
			return models.UserPhoneCode{}, false, nil
		}
		return models.UserPhoneCode{}, false, err
	}
	var phoneCode models.UserPhoneCode
	if err = json.Unmarshal([]byte(data), &phoneCode); err != nil {
		return models.UserPhoneCode{}, true, err
	}
	return phoneCode, true, nil
}

func DeleteUserPhoneCode(accountId string) error {
	return GetCache().Del(NewUserPhoneCodeKey(accountId))
}
