package application

import (
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"context"
)

type UserService struct {
}

func (service *UserService) UserDaoRegister(ctx context.Context, req *user.ReqUserRegister) (*common.CommonResponse, error) {
	return &common.CommonResponse{}, nil
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
