package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CoinLogHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	pageReq, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &user.ReqUserCoinLogQueryDao{
		AccountId: userMeta.AccountId,
		Page:      pageReq.Page,
		PageSize:  pageReq.PageSize,
	}
	var resp *user.RspUserCoinLogQueryDao
	resp, err = userClient.UserCoinLogQueryDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, cl := range resp.List {
		var coinLog = rsp.UserCoinLog{
			Id:        cl.Id,
			Coin:      cl.Coin,
			CoinType:  constant.CoinType(cl.CoinType),
			Describe:  cl.Describe,
			Timestamp: cl.Timestamp,
		}
		list = append(list, coinLog)
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

func DeleteCoinLogHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	idsStr := strings.Split(c.Query("ids"), ",")
	if len(idsStr) == 0 {
		err := fmt.Errorf("id list can't null")
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var ids = make([]int64, 0)
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Logger.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
			return
		}
		ids = append(ids, int64(id))
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &user.ReqUserCoinLogDeleteDao{
		AccountId: userMeta.AccountId,
		Ids:       ids,
	}
	_, err = userClient.UserCoinLogDeleteDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func CoinRankHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *user.RspUserCoinRankQueryDao
	resp, err = userClient.UserCoinRankQueryDao(context.Background(), &user.ReqUserCoinRankQueryDao{AccountId: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var response = rsp.UserCoinRank{
		User: rsp.User{
			AccountID: userMeta.AccountId,
			Username:  userMeta.Username,
		},
		Rank:          resp.Rank,
		Coin:          resp.Coin,
		SameCoinUsers: resp.SameCoinUsers,
	}
	handle.SuccessResp(c, "", response)
}

func CoinRankBoardHandle(c *gin.Context) {
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *user.RspUserCoinRankBoardQueryDao
	resp, err = userClient.UserCoinRankBoardQueryDao(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, ucr := range resp.List {
		ids = append(ids, ucr.AccountId)
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var idsMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		idsMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
		}
	}
	var list = make([]rsp.UserCoinRankBoard, 0)
	for _, ucr := range resp.List {
		var userCoinRank = rsp.UserCoinRankBoard{
			User:            idsMap[ucr.AccountId],
			Rank:            ucr.Rank,
			Coin:            ucr.Coin,
			UpdateTimestamp: ucr.UpdateTimestamp,
		}
		list = append(list, userCoinRank)
	}
	var response = rsp.UserCoinRankBoardResp{
		List:               list,
		StatisticTimestamp: time.Now().Unix(),
	}
	handle.SuccessResp(c, "", response)
}
