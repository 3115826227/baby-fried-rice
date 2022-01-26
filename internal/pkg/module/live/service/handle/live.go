package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/live"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/live/grpc"
	"baby-fried-rice/internal/pkg/module/live/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 直播房间列表查询
func LiveRoomHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var client live.DaoLiveClient
	if client, err = grpc.GetLiveClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var req = live.ReqLiveRoomQueryDao{
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
	}
	var resp *live.RspLiveRoomQueryDao
	if resp, err = client.LiveRoomQueryDao(context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, item := range resp.List {
		ids = append(ids, item.Origin)
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		userMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var list = make([]interface{}, 0)
	for _, item := range resp.List {
		list = append(list, rsp.LiveRoom{
			LiveRoomId: item.LiveRoomId,
			Origin:     userMap[item.Origin],
			Status:     item.Status,
			UserTotal:  item.UserTotal,
		})
	}
	handle.SuccessListResp(c, "", list, 0, reqPage.Page, reqPage.PageSize)
}

// 直播房间详情查询
func LiveRoomDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	client, err := grpc.GetLiveClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var req = live.ReqLiveRoomDetailQueryDao{
		LiveRoomId: c.Query("live_room_id"),
		AccountId:  userMeta.AccountId,
	}
	var resp *live.RspLiveRoomDetailQueryDao
	if resp, err = client.LiveRoomDetailQueryDao(context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{resp.Origin}}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	if len(userResp.Users) != 1 {
		err = fmt.Errorf("get user dao by id %v failed", resp.Origin)
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var u = userResp.Users[0]
	var originUser = rsp.User{
		AccountID:  u.Id,
		Username:   u.Username,
		HeadImgUrl: u.HeadImgUrl,
		IsOfficial: u.IsOfficial,
	}
	var response = rsp.LiveRoomDetailResp{
		LiveRoom: rsp.LiveRoom{
			LiveRoomId: resp.LiveRoomId,
			Origin:     originUser,
			Status:     resp.Status,
			UserTotal:  resp.UserTotal,
		},
		OnlineTime: resp.OnlineTime,
	}
	handle.SuccessResp(c, "", response)
}

// 直播房间主播用户列表查询
func LiveRoomUserHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var client live.DaoLiveClient
	if client, err = grpc.GetLiveClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var req = live.ReqLiveRoomUserQueryDao{
		Origin:     userMeta.AccountId,
		LiveRoomId: c.Query("live_room_id"),
		Page:       reqPage.Page,
		PageSize:   reqPage.PageSize,
	}
	var resp *live.RspLiveRoomUserQueryDao
	if resp, err = client.LiveRoomUserQueryDao(context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, u := range resp.Users {
		ids = append(ids, u)
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		userMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var list = make([]rsp.User, 0)
	for _, u := range resp.Users {
		list = append(list, userMap[u])
	}
	var response = rsp.LiveRoomUserResp{
		List:     list,
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
	}
	handle.SuccessResp(c, "", response)
}

// 直播房间主播更新操作
func LiveRoomOriginUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqUpdateOriginLiveRoom
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var client live.DaoLiveClient
	if client, err = grpc.GetLiveClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *live.RspLiveRoomDetailQueryDao
	if resp, err = client.LiveRoomStatusUpdateDao(context.Background(), &live.ReqLiveRoomStatusUpdateDao{
		Origin: userMeta.AccountId,
		Status: req.Status,
	}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var response = rsp.UpdateOriginLiveRoomResp{
		LiveRoom: rsp.LiveRoom{
			LiveRoomId: resp.LiveRoomId,
			Origin:     userMeta.GetUser(),
			Status:     resp.Status,
		},
	}
	switch req.Status {
	case live.LiveRoomStatus_Online:
		var swapSdp string
		swapSdp, err = CreateSession(req.Sdp, resp.LiveRoomId, userMeta.AccountId, true)
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
			return
		}
		response.SwapSdp = swapSdp
	}
	handle.SuccessResp(c, "", response)
}

// 直播房间用户更新操作
func LiveRoomUserOptUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqUpdateUserLiveRoomOpt
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var client live.DaoLiveClient
	if client, err = grpc.GetLiveClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var liveReq = live.ReqLiveRoomUserOptAddDao{
		LiveRoomId: req.LiveRoomId,
		AccountId:  userMeta.AccountId,
		Opt:        req.Opt,
	}
	var resp *live.RspLiveRoomUserOptAddDao
	if resp, err = client.LiveRoomUserOptAddDao(context.Background(), &liveReq); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var response rsp.UpdateUserLiveRoomOptResp
	switch req.Opt {
	case live.LiveRoomUserOptType_EnterOptType:
		var remoteSwapSdp string
		if remoteSwapSdp, err = JoinSession(req.RemoteSdp, req.LiveRoomId, resp.Origin, true); err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
			return
		}
		response.RemoteSwapSdp = remoteSwapSdp
	}
	handle.SuccessResp(c, "", response)
}

// 直播房间消息查询
func LiveRoomUserMessageHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var client live.DaoLiveClient
	if client, err = grpc.GetLiveClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var req = live.ReqLiveRoomMessageQueryDao{
		AccountId:  userMeta.AccountId,
		LiveRoomId: c.Query("live_room_id"),
		Page:       reqPage.Page,
		PageSize:   reqPage.PageSize,
	}
	var resp *live.RspLiveRoomMessageQueryDao
	if resp, err = client.LiveRoomMessageQueryDao(context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, msg := range resp.List {
		ids = append(ids, msg.Send)
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		userMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var list = make([]interface{}, 0)
	for _, msg := range resp.List {
		list = append(list, rsp.LiveRoomMessage{
			MessageId:     msg.MessageId,
			MessageType:   msg.MessageType,
			Send:          userMap[msg.Send],
			Content:       msg.Content,
			SendTimestamp: msg.SendTimestamp,
		})
	}
	handle.SuccessListResp(c, "", list, 0, reqPage.Page, reqPage.PageSize)
}
