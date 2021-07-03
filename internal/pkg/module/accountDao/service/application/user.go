package application

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/config"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type UserService struct {
}

func (service *UserService) UserDaoById(ctx context.Context, req *user.ReqUserDaoById) (resp *user.RspUserDaoById, err error) {
	var details []tables.AccountUserDetail
	details, err = query.GetUsers(req.Ids)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var users = make([]*user.UserDao, 0)
	for _, detail := range details {
		users = append(users, &user.UserDao{
			Id:         detail.AccountID,
			Username:   detail.Username,
			HeadImgUrl: detail.HeadImgUrl,
		})
	}
	resp = &user.RspUserDaoById{Users: users}
	return
}

func (service *UserService) UserDaoRegister(ctx context.Context, req *user.ReqUserRegister) (empty *emptypb.Empty, err error) {
	if query.IsDuplicateLoginNameByUser(req.Login.LoginName) {
		log.Logger.Error(fmt.Sprintf("login name %v is duplication", req.Login.LoginName))
		return
	}
	accountID := handle.GenerateSerialNumber()

	var now = time.Now()
	var accountUser tables.AccountUser
	accountUser.ID = handle.GenerateID()
	accountUser.AccountId = accountID
	accountUser.LoginName = req.Login.LoginName
	accountUser.Password = req.Login.Password
	accountUser.EncodeType = config.DefaultUserEncryMd5
	accountUser.CreatedAt = now
	accountUser.UpdatedAt = now

	var detail tables.AccountUserDetail
	detail.ID = accountUser.ID
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

	var beans = make([]interface{}, 0)
	beans = append(beans, &accountUser)
	beans = append(beans, &detail)

	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
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
	if err = cache.AddUserDetail(*detail); err != nil {
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
			Age:       detail.Age,
			Phone:     detail.Phone,
		},
	}
	return
}

func (service *UserService) UserDaoDetail(ctx context.Context, req *user.ReqDaoUserDetail) (resp *user.RspDaoUserDetail, err error) {
	var detail tables.AccountUserDetail
	detail, err = query.GetUserDetail(req.AccountId)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspDaoUserDetail{
		Detail: &user.DaoUserDetail{
			AccountId:  detail.AccountID,
			HeadImgUrl: detail.HeadImgUrl,
			Username:   detail.Username,
			SchoolId:   detail.SchoolId,
			Gender:     detail.Gender,
			Age:        detail.Age,
			Phone:      detail.Phone,
			Describe:   detail.Describe,
		},
	}
	return
}

func (service *UserService) UserDaoDetailUpdate(ctx context.Context, req *user.ReqDaoUserDetailUpdate) (empty *emptypb.Empty, err error) {
	var detail = tables.AccountUserDetail{
		Username:   req.Detail.Username,
		SchoolId:   req.Detail.SchoolId,
		Gender:     req.Detail.Gender,
		Age:        req.Detail.Age,
		HeadImgUrl: req.Detail.HeadImgUrl,
		Phone:      req.Detail.Phone,
		Describe:   req.Detail.Describe,
	}
	if err = db.GetDB().GetDB().Where("account_id = ?", req.Detail.AccountId).Updates(&detail).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var newDetail tables.AccountUserDetail
	if err = db.GetDB().GetObject(map[string]interface{}{"account_id": req.Detail.AccountId}, &newDetail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = cache.AddUserDetail(newDetail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *UserService) UserDaoPwdUpdate(ctx context.Context, req *user.ReqDaoUserPwdUpdate) (empty *emptypb.Empty, err error) {
	var accountUser tables.AccountUser
	if err = db.GetDB().GetObject(map[string]interface{}{"account_id": req.AccountId}, &accountUser); err != nil {
		return
	}
	accountUser.Password = req.NewPassword
	accountUser.UpdatedAt = time.Now()
	if err = db.GetDB().GetDB().Where("account_id = ? and password = ?", req.AccountId, req.Password).Updates(&accountUser).Error; err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *UserService) UserDaoAll(ctx context.Context, empty *emptypb.Empty) (resp *user.RspUserDaoAll, err error) {
	var ids []string
	if ids, err = query.GetAll(); err != nil {
		return
	}
	resp = &user.RspUserDaoAll{AccountIds: ids}
	return
}
