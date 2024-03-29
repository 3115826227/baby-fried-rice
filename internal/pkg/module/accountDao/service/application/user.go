package application

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	Errors "baby-fried-rice/internal/pkg/kit/errors"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/kit/utils"
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserService struct {
}

func (service *UserService) UserDaoById(ctx context.Context, req *user.ReqUserDaoById) (resp *user.RspUserDaoById, err error) {
	var details []tables.AccountUserDetail
	details, err = query.GetUserDetails(req.Ids)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var users = make([]*user.UserDao, 0)
	for _, detail := range details {
		var phoneVerify = false
		if detail.Phone != "" {
			phoneVerify = true
		}
		users = append(users, &user.UserDao{
			Id:          detail.AccountID,
			Username:    detail.Username,
			HeadImgUrl:  detail.HeadImgUrl,
			IsOfficial:  detail.IsOfficial,
			PhoneVerify: phoneVerify,
		})
	}
	resp = &user.RspUserDaoById{Users: users}
	return
}

func (service *UserService) UserDaoLoginNameExist(ctx context.Context, req *user.ReqUserDaoLoginNameExist) (resp *user.RspUserDaoLoginNameExist, err error) {
	var exist bool
	exist, err = query.IsDuplicateLoginNameByUser(req.LoginName)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("login name %v is duplication", req.LoginName))
		err = Errors.NewCommonError(constant.CodeLoginNameExist)
		return
	}
	resp = &user.RspUserDaoLoginNameExist{
		Exist: exist,
	}
	return
}

func (service *UserService) UserDaoRegister(ctx context.Context, req *user.ReqUserRegister) (empty *emptypb.Empty, err error) {
	var exist bool
	exist, err = query.IsDuplicateLoginNameByUser(req.Login.LoginName)
	if err != nil {
		log.Logger.Error(err.Error())
		err = Errors.NewCommonErr(constant.CodeInvalidParams, err)
		return nil, Errors.ConvertEdgeXErrToRpc(err)
	}
	if exist {
		log.Logger.Error(fmt.Sprintf("login name %v is duplication", req.Login.LoginName))
		err = Errors.NewCommonErr(constant.CodeLoginNameExist, err)
		return nil, Errors.ConvertEdgeXErrToRpc(err)
	}
	accountID := handle.GenerateSerialNumber()
	for {
		if !query.IsDuplicateAccountID(accountID) {
			break
		}
	}

	var now = time.Now()
	var accountUser tables.AccountUser
	accountUser.ID = handle.GenerateID()
	accountUser.AccountId = accountID
	accountUser.LoginName = req.Login.LoginName
	accountUser.Password = handle.EncodePassword(accountID, req.Login.Password)
	accountUser.EncodeType = constant.DefaultUserEncryMd5
	accountUser.CreatedAt = now
	accountUser.UpdatedAt = now

	var detail tables.AccountUserDetail
	detail.ID = accountUser.ID

	detail.AccountID = accountID
	detail.Username = req.Username
	detail.CreatedAt = now
	detail.UpdatedAt = now

	var coin tables.AccountUserCoin
	coin.AccountID = accountID
	/*
		新用户注册，送500积分
	*/
	coin.Coin = constant.UserRegisterCoin
	coin.CoinTotal = coin.Coin
	coin.UpdateTimestamp = now.Unix()

	var coinLog = tables.AccountUserCoinLog{
		AccountID: accountID,
		Coin:      coin.Coin,
		CoinType:  constant.UserRegisterCoinType,
		Timestamp: now.Unix(),
	}

	var beans = make([]interface{}, 0)
	beans = append(beans, &accountUser)
	beans = append(beans, &detail)
	beans = append(beans, &coin)
	beans = append(beans, &coinLog)

	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return nil, Errors.ConvertEdgeXErrToRpc(err)
	}
	// 将用户信息和用户积分信息写入缓存
	if err = cache.SetUserDetail(detail); err != nil {
		log.Logger.Error(err.Error())
		return nil, Errors.ConvertEdgeXErrToRpc(err)
	}
	if err = cache.SetUserCoin(coin); err != nil {
		log.Logger.Error(err.Error())
		return nil, Errors.ConvertEdgeXErrToRpc(err)
	}
	empty = new(emptypb.Empty)
	return
}

func (service *UserService) UserDaoLogin(ctx context.Context, req *user.ReqPasswordLogin) (resp *user.RspDaoUserLogin, err error) {
	var loginUser tables.AccountUser
	loginUser, err = query.GetUserByLogin(req.LoginName)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if loginUser.Password != handle.EncodePassword(loginUser.AccountId, req.Password) {
		err = errors.New("password is invalid")
		log.Logger.Error(err.Error())
		return
	}
	var detail tables.AccountUserDetail
	if err = db.GetDB().GetObject(map[string]interface{}{"account_id": loginUser.AccountId}, &detail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = cache.SetUserDetail(detail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 写入日志
	go func() {
		if !loginUser.Cancel && !loginUser.Freeze {
			var loginLog = tables.AccountUserLoginLog{
				AccountId: detail.AccountID,
				IP:        req.Ip,
				LoginTime: time.Now(),
			}
			if err = db.GetDB().CreateObject(&loginLog); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	}()
	resp = &user.RspDaoUserLogin{
		User: &user.RspDaoUser{
			AccountId: detail.AccountID,
			LoginName: loginUser.LoginName,
			Username:  detail.Username,
			SchoolId:  detail.SchoolId,
			Gender:    detail.Gender,
			Age:       detail.Age,
			Phone:     detail.Phone,
			Freeze:    loginUser.Freeze,
			Cancel:    loginUser.Cancel,
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
	var coin tables.AccountUserCoin
	coin, err = query.GetUserCoin(req.AccountId)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Logger.Error(err.Error())
			return
		}
		err = nil
		coin = tables.AccountUserCoin{
			AccountID: detail.AccountID,
		}
	}
	resp = &user.RspDaoUserDetail{
		Detail: &user.DaoUserDetail{
			AccountId:  detail.AccountID,
			HeadImgUrl: detail.HeadImgUrl,
			Username:   detail.Username,
			SchoolId:   detail.SchoolId,
			Gender:     detail.Gender,
			Age:        detail.Age,
			Describe:   detail.Describe,
			Coin:       coin.Coin,
			IsOfficial: detail.IsOfficial,
		},
	}
	if detail.Phone != "" {
		var decodePhoneBytes []byte
		decodePhoneBytes, err = utils.Base64Decode(detail.Phone)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var phoneBytes []byte
		phoneBytes, err = utils.GcmDecrypt(decodePhoneBytes, utils.EncryptKey(detail.AccountID))
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		resp.Detail.Phone = string(phoneBytes)
	}
	return
}

func (service *UserService) UserDaoDetailUpdate(ctx context.Context, req *user.ReqDaoUserDetailUpdate) (empty *emptypb.Empty, err error) {
	var encodePhone string
	if req.Detail.Phone != "" {
		// 手机号验证
		var encodePhoneBytes []byte
		encodePhoneBytes, err = utils.GcmEncrypt([]byte(req.Detail.Phone), utils.EncryptKey(req.Detail.AccountId))
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		encodePhone = utils.Base64Encode(encodePhoneBytes)
	}
	var now = time.Now()
	var detail = tables.AccountUserDetail{
		AccountID:  req.Detail.AccountId,
		Username:   req.Detail.Username,
		SchoolId:   req.Detail.SchoolId,
		Gender:     req.Detail.Gender,
		Age:        req.Detail.Age,
		HeadImgUrl: req.Detail.HeadImgUrl,
		Describe:   req.Detail.Describe,
		Phone:      encodePhone,
	}
	detail.UpdatedAt = now
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	if encodePhone != "" {
		var phone = tables.AccountUserPhone{
			Phone: req.Detail.Phone,
		}
		phone.CreatedAt, phone.UpdatedAt = now, now
		if err = tx.Create(&phone).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	// 更新数据库的用户修改信息
	if err = tx.Where("account_id = ?", req.Detail.AccountId).Updates(&detail).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 删除缓存信息
	if err = cache.DeleteUserDetail(detail.AccountID); err != nil {
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
	if accountUser.Password != handle.EncodePassword(accountUser.AccountId, req.Password) {
		err = errors.New("old password is invalid")
		log.Logger.Error(err.Error())
		return
	}
	accountUser.Password = handle.EncodePassword(accountUser.AccountId, req.NewPassword)
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

// 查询手机号是否被认证
func (service *UserService) UserDaoPhoneVerify(ctx context.Context, req *user.ReqUserDaoPhoneVerify) (resp *user.RspUserDaoPhoneVerify, err error) {
	var exist bool
	if exist, err = query.GetUserPhone(req.Phone); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspUserDaoPhoneVerify{
		Verify: exist,
	}
	return
}

// 用户积分信息查询
func (service *UserService) UserCoinDao(ctx context.Context, req *user.ReqUserCoinDao) (resp *user.RspUserCoinDao, err error) {
	var coin tables.AccountUserCoin
	coin, err = query.GetUserCoin(req.AccountId)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspUserCoinDao{
		AccountId: coin.AccountID,
		Coin:      coin.Coin,
	}
	return
}

// 积分变动，积分日志添加
func (service *UserService) UserCoinLogAddDao(ctx context.Context, req *user.ReqUserCoinLogAddDao) (empty *emptypb.Empty, err error) {
	tx := db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	// 获取用户积分
	var detail tables.AccountUserCoin
	if err = db.GetDB().GetObject(map[string]interface{}{"account_id": req.AccountId}, &detail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 计算新的总积分和总获取积分
	var newCoin, newCoinTotal = req.Coin + detail.Coin, detail.CoinTotal
	if newCoin < 0 {
		err = errors.New("user coin insufficient")
		log.Logger.Error(err.Error())
		return
	}
	if req.Coin > 0 {
		newCoinTotal += req.Coin
	}
	var now = time.Now()
	// 更新新的总积分和总获取积分
	if err = tx.Model(&tables.AccountUserCoin{}).Where("account_id = ?",
		req.AccountId).Updates(map[string]interface{}{
		"coin":             newCoin,
		"coin_total":       newCoinTotal,
		"update_timestamp": now.Unix(),
	}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 删除积分缓存
	if err = cache.DeleteUserCoin(req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 添加积分记录
	var coinLog = tables.AccountUserCoinLog{
		AccountID: req.AccountId,
		Coin:      req.Coin,
		CoinType:  constant.CoinType(req.CoinType),
		Timestamp: now.Unix(),
	}
	if err = tx.Create(&coinLog).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 用户积分日志查询
func (service *UserService) UserCoinLogQueryDao(ctx context.Context, req *user.ReqUserCoinLogQueryDao) (resp *user.RspUserCoinLogQueryDao, err error) {
	var pageReq = requests.PageCommonReq{
		PageSize: req.PageSize,
		Page:     req.Page,
	}
	var logs []tables.AccountUserCoinLog
	var total int64
	if logs, total, err = query.GetUserCoinLog(req.AccountId, pageReq); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*user.UserCoinLogQueryDao, 0)
	for _, l := range logs {
		coinLogDao := &user.UserCoinLogQueryDao{
			Id:        l.ID,
			AccountId: l.AccountID,
			Coin:      l.Coin,
			CoinType:  int64(l.CoinType),
			Describe:  constant.CoinTypeDescribeMap[l.CoinType],
			Timestamp: l.Timestamp,
		}
		list = append(list, coinLogDao)
	}
	resp = &user.RspUserCoinLogQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}

// 用户积分日志删除
func (service *UserService) UserCoinLogDeleteDao(ctx context.Context, req *user.ReqUserCoinLogDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("account_id = ? and id in (?)",
		req.AccountId, req.Ids).Delete(&tables.AccountUserCoinLog{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 用户积分排名查询
func (service *UserService) UserCoinRankQueryDao(ctx context.Context, req *user.ReqUserCoinRankQueryDao) (resp *user.RspUserCoinRankQueryDao, err error) {
	resp = new(user.RspUserCoinRankQueryDao)
	resp.AccountId = req.AccountId
	if resp.Rank, resp.SameCoinUsers, resp.Coin, err = query.GetUserCoinRank(req.AccountId); err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		log.Logger.Error(err.Error())
		return
	}
	return
}

// 用户积分排行榜查询
func (service *UserService) UserCoinRankBoardQueryDao(ctx context.Context, empty *emptypb.Empty) (resp *user.RspUserCoinRankBoardQueryDao, err error) {
	var usersMap []map[string]interface{}
	if usersMap, err = query.GetUserCoinRankBoard(constant.UserCoinBoardTotal); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*user.UserCoinRankBoardQueryDao, 0)
	for _, userMap := range usersMap {
		var u = &user.UserCoinRankBoardQueryDao{
			AccountId:       userMap["account_id"].(string),
			Rank:            userMap["rank"].(int64),
			Coin:            userMap["coin"].(int64),
			UpdateTimestamp: userMap["timestamp"].(int64),
		}
		list = append(list, u)
	}
	resp = &user.RspUserCoinRankBoardQueryDao{
		List:               list,
		StatisticTimestamp: time.Now().Unix(),
	}
	return
}

// 用户签到
func (service *UserService) UserSignInDao(ctx context.Context, req *user.ReqUserSignInDao) (resp *user.RspUserSignInDao, err error) {
	var (
		latestSignInLog tables.AccountUserSignInLog       // 最近一条签到日志
		signInCoin      constant.RewardCoinBySignedInType // 推算出来的签到积分奖励
	)
	// 查出最近一条签到日志
	if latestSignInLog, err = query.GetUserLatestSignIn(req.AccountId); err != nil {
		// 没有找到记录 按照第一次签到计算
		if err == gorm.ErrRecordNotFound {
			signInCoin = constant.RewardCoinBySignedInContinuedOneDay
		} else {
			// 查询数据库出错，直接返回
			log.Logger.Error(err.Error())
			return
		}
	} else {
		// 判断这条日志是当日签到/昨日签到/昨日之前签到
		var t = time.Now()
		var todayStartTimestamp = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		var yesterdayStartTimestamp = time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
		if latestSignInLog.Timestamp > todayStartTimestamp.Unix() {
			// 当日已经签到，直接返回结果
			resp = &user.RspUserSignInDao{
				AccountId: req.AccountId,
				Ok:        false,
				Describe:  fmt.Sprintf("今日您的签到已完成，不需要再次签到，您已连续签到%v天", constant.RewardCoinBySignInContinueDayMap[latestSignInLog.Coin]),
				Coin:      0,
			}
			return
		} else {
			if latestSignInLog.Timestamp > yesterdayStartTimestamp.Unix() {
				// 昨日已签到，根据昨日签到的积分奖励进行推算
				signInCoin = constant.RewardCoinBySignedInNextMap[latestSignInLog.Coin]
			} else {
				// 昨日未签到，按照第一天签到计算
				signInCoin = constant.RewardCoinBySignedInContinuedOneDay
			}
		}
	}
	// 查询用户现有积分，计算出新的需要更新积分值
	var userCoin tables.AccountUserCoin
	if userCoin, err = query.GetUserCoin(req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var now = time.Now()
	var newCoin, newCoinTotal = userCoin.Coin + int64(signInCoin), userCoin.CoinTotal + int64(signInCoin)
	// 更新积分
	if err = tx.Model(&tables.AccountUserCoin{}).Where("account_id = ?", req.AccountId).Updates(map[string]interface{}{
		"coin":             newCoin,
		"coin_total":       newCoinTotal,
		"update_timestamp": now.Unix(),
	}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 添加积分日志
	var coinLog = tables.AccountUserCoinLog{
		AccountID: req.AccountId,
		Coin:      int64(signInCoin),
		CoinType:  constant.SignInCoinType,
		Timestamp: now.Unix(),
	}
	if err = tx.Create(&coinLog).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 添加签到日志
	var signInLog = tables.AccountUserSignInLog{
		AccountId:  req.AccountId,
		Coin:       signInCoin,
		SignInType: constant.NormalSignInType,
		Timestamp:  now.Unix(),
	}
	if err = tx.Create(&signInLog).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspUserSignInDao{
		AccountId: req.AccountId,
		Ok:        true,
		Describe:  fmt.Sprintf("连续签到第%v天，奖励%v积分", constant.RewardCoinBySignInContinueDayMap[signInCoin], signInCoin),
		Coin:      int64(signInCoin),
	}
	// 删除缓存中的用户积分
	if err = cache.DeleteUserCoin(req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 删除缓存中的最近更新日志
	if err = cache.DeleteUserSignInLatestLog(req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 更新缓存中的用户积分排名
	if err = cache.SetUserCoinRank(coinLog, int64(signInCoin)); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}

// 用户签到记录查询
func (service *UserService) UserSignInLogQueryDao(ctx context.Context, req *user.ReqUserSignInLogQueryDao) (resp *user.RspUserSignInLogQueryDao, err error) {
	var startTimestamp = time.Now().AddDate(int(-req.Year), int(-req.Month), int(-req.Day)).Unix()
	var signInLogs []tables.AccountUserSignInLog
	if err = db.GetDB().GetDB().Where("account_id = ? and timestamp > ?",
		req.AccountId, startTimestamp).Order("timestamp desc").Find(&signInLogs).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*user.UserSignInLogQueryDao, 0)
	for _, sl := range signInLogs {
		var signInLog = &user.UserSignInLogQueryDao{
			AccountId:  sl.AccountId,
			Timestamp:  sl.Timestamp,
			Coin:       int64(sl.Coin),
			SignInType: int64(sl.SignInType),
		}
		list = append(list, signInLog)
	}
	resp = &user.RspUserSignInLogQueryDao{List: list}
	return
}

func (service *UserService) UserCommunicationAddDao(ctx context.Context, req *user.ReqUserCommunicationAddDao) (resp *user.RspUserCommunicationAddDao, err error) {
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var communication = tables.Communication{
		Title:             req.Title,
		Origin:            req.AccountId,
		CommunicationType: req.CommunicationType,
	}
	var now = time.Now()
	communication.CreatedAt, communication.UpdatedAt = now, now
	if err = tx.Create(&communication).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var detail = tables.CommunicationDetail{
		Id:      communication.ID,
		Content: req.Content,
		Images:  req.Images,
	}
	if err = tx.Create(&detail).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspUserCommunicationAddDao{Id: detail.Id}
	return
}

func (service *UserService) UserCommunicationQueryDao(ctx context.Context, req *user.ReqUserCommunicationQueryDao) (resp *user.RspUserCommunicationQueryDao, err error) {
	var params = query.CommunicationQueryParams{
		CommunicationType: req.CommunicationType,
		Page:              req.Page,
		PageSize:          req.PageSize,
	}
	var communications []tables.Communication
	var total int64
	if communications, total, err = query.GetCommunication(params); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*user.UserCommunicationDao, 0)
	for _, communication := range communications {
		list = append(list, &user.UserCommunicationDao{
			Id:                communication.ID,
			Origin:            communication.Origin,
			Title:             communication.Title,
			CommunicationType: communication.CommunicationType,
			CreateTimestamp:   communication.CreatedAt.Unix(),
			UpdateTimestamp:   communication.UpdatedAt.Unix(),
			Reply:             communication.Reply,
		})
	}
	resp = &user.RspUserCommunicationQueryDao{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return
}

func (service *UserService) UserCommunicationDetailQueryDao(ctx context.Context, req *user.ReqUserCommunicationDetailQueryDao) (resp *user.RspUserCommunicationDetailQueryDao, err error) {
	var communication tables.Communication
	var detail tables.CommunicationDetail
	communication, detail, err = query.GetCommunicationDetail(req.Id, req.Origin)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &user.RspUserCommunicationDetailQueryDao{
		Communication: &user.UserCommunicationDao{
			Id:                communication.ID,
			Origin:            communication.Origin,
			Title:             communication.Title,
			CommunicationType: communication.CommunicationType,
			CreateTimestamp:   communication.CreatedAt.Unix(),
			UpdateTimestamp:   communication.UpdatedAt.Unix(),
			Reply:             communication.Reply,
		},
		Content:        detail.Content,
		ReplyContent:   detail.ReplyContent,
		ReplyTimestamp: detail.ReplyTimestamp,
	}
	if detail.Images != "" {
		resp.Images = strings.Split(detail.Images, ",")
	}
	return
}

func (service *UserService) UserCommunicationDeleteDao(ctx context.Context, req *user.ReqUserCommunicationDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Delete(&tables.Communication{}, "id = ? and origin = ?",
		req.Id, req.Origin).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *UserService) IteratorVersionQueryDao(ctx context.Context, empty *emptypb.Empty) (resp *user.RspIteratorVersionQueryDao, err error) {
	var versions []tables.IterativeVersion
	versions, err = query.GetIteratorVersion()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*user.IteratorVersionDao, 0)
	for _, v := range versions {
		list = append(list, &user.IteratorVersionDao{
			Version:   v.Version,
			Content:   v.Content,
			Timestamp: v.UpdateTimestamp,
		})
	}
	resp = &user.RspIteratorVersionQueryDao{List: list}
	return
}

// 用户举报
func (service *UserService) UserTipAddDao(ctx context.Context, req *user.ReqUserTipAddDao) (empty *emptypb.Empty, err error) {
	var now = time.Now()
	var tip = tables.AccountUserTipLog{
		ReportAccountId:   req.ReportAccountId,
		ReportedAccountId: req.ReportedAccountId,
		Describe:          req.Describe,
	}
	tip.CreatedAt, tip.UpdatedAt = now, now
	if err = db.GetDB().CreateObject(&tip); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}
