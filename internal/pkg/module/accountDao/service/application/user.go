package application

import (
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/user"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/accountDao/config"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"context"
	"fmt"
	"time"
)

type UserService struct {
}

func (service *UserService) UserDaoRegister(ctx context.Context, req *user.ReqUserRegister) (resp *common.CommonResponse, err error) {
	if query.IsDuplicateLoginNameByUser(req.Login.LoginName) {
		log.Logger.Error(fmt.Sprintf("login name %v is duplication", req.Login.LoginName))
		return
	}

	var now = time.Now()
	var accountUser tables.AccountUser
	accountUser.ID = handle.GenerateID()
	accountUser.LoginName = req.Login.LoginName
	accountUser.Password = req.Login.Password
	accountUser.EncodeType = config.DefaultUserEncryMd5
	accountUser.CreatedAt = now
	accountUser.UpdatedAt = now

	var detail tables.AccountUserDetail
	detail.ID = accountUser.ID
	accountID := handle.GenerateSerialNumber()
	for {
		if !query.IsDuplicateAccountID(accountID) {
			break
		}
	}

	detail.AccountID = accountID
	detail.Username = req.Username
	detail.Gender = req.Gender
	detail.Phone = req.Phone
	detail.CreatedAt = now
	detail.UpdatedAt = now

	var userDetail tables.UserDetail
	userDetail.UserId = detail.ID
	userDetail.AccountId = detail.AccountID
	userDetail.Username = detail.Username

	var beans = make([]interface{}, 0)
	beans = append(beans, &accountUser)
	beans = append(beans, &detail)
	beans = append(beans, &userDetail)

	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}

func (service *UserService) UserDaoLogin(ctx context.Context, req *user.ReqPasswordLogin) (resp *user.RspDaoUserLogin, err error) {
	loginUser, err := query.GetUserByLogin(req.LoginName, req.Password)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var detail = new(tables.AccountUserDetail)
	detail.ID = loginUser.ID
	if err = db.GetDB().GetObject(nil, detail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspDaoUserLogin{
		User: &user.RspDaoUser{
			AccountId: detail.AccountID,
			LoginName: loginUser.LoginName,
			Username:  detail.Username,
			SchoolId:  detail.SchoolId,
			Gender:    detail.Gender,
			Age:       int64(detail.Age),
			Phone:     detail.Phone,
		},
	}
	return
}
