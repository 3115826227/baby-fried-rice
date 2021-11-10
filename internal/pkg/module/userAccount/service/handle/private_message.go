package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 发送私信
func SendPrivateMessageHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var pmReq requests.UserSendPrivateMessageReq
	pmReq.SendId = userMeta.AccountId
	if err := c.ShouldBind(&pmReq); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	pmClient, err := grpc.GetPrivateMessageClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var now = time.Now().Unix()
	reqAdd := privatemessage.ReqPrivateMessageAddDao{
		SendId:          pmReq.SendId,
		ReceiveId:       pmReq.ReceiveId,
		MessageType:     pmReq.MessageType,
		MessageSendType: int32(pmReq.MessageSendType),
		Title:           pmReq.MessageTitle,
		Content:         pmReq.MessageContent,
		CreateTimestamp: now,
	}
	var resp *privatemessage.RspPrivateMessageAddDao
	resp, err = pmClient.PrivateMessageAddDao(context.Background(), &reqAdd)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		reqDetail := &privatemessage.ReqPrivateMessageDetailDao{
			AccountId: userMeta.AccountId,
			Id:        resp.Id,
		}
		var respDetail *privatemessage.RspPrivateMessageDetailDao
		respDetail, err = pmClient.PrivateMessageDetailDao(context.Background(), reqDetail)
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		var pm = respDetail.PrivateMessage
		var detail = &user.UserDao{
			Id:         userMeta.AccountId,
			Username:   userMeta.Username,
			IsOfficial: userMeta.IsOfficial,
		}
		var notify = models.WSMessageNotify{
			WSMessageNotifyType: constant.PrivateMessageNotify,
			Receive:             pmReq.ReceiveId,
			WSMessage: models.WSMessage{
				Send: userMeta.GetUser(),
				PrivateMessage: rsp.UserPrivateMessageDetailResp{
					UserPrivateMessage: rsp.PrivateMessagePbConvertToRsp(pm, detail),
					Content:            respDetail.Content,
				},
			},
			Timestamp: now,
		}
		if err = mq.Send(config.GetConfig().MessageQueue.PublishTopics.WebsocketNotify, notify.ToString()); err != nil {
			log.Logger.Error(err.Error())
		}
	}()
	handle.SuccessResp(c, "", nil)
}

// 查询私信列表
func PrivateMessagesHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sendId := c.Query("send_id")
	pageReq, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var pmClient privatemessage.DaoPrivateMessageClient
	pmClient, err = grpc.GetPrivateMessageClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	reqQuery := privatemessage.ReqPrivateMessageQueryDao{
		Page:      pageReq.Page,
		PageSize:  pageReq.PageSize,
		SendId:    sendId,
		AccountId: userMeta.AccountId,
	}
	var resp *privatemessage.RspPrivateMessageQueryDao
	resp, err = pmClient.PrivateMessageQueryDao(context.Background(), &reqQuery)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, pm := range resp.List {
		ids = append(ids, pm.SendId)
	}
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var idsMap = make(map[string]*user.UserDao)
	for _, u := range userResp.Users {
		idsMap[u.Id] = &user.UserDao{
			Id:         u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var list = make([]interface{}, 0)
	for _, pm := range resp.List {
		list = append(list, rsp.PrivateMessagePbConvertToRsp(pm, idsMap[pm.SendId]))
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

// 更新私信读取状态
func UpdatePrivateMessageStatusHandle(c *gin.Context) {
	var userMeta = handle.GetUserMeta(c)
	var upm requests.UpdatePrivateMessageStatusReq
	if err := c.ShouldBind(&upm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	if len(upm.MessageIds) != 0 {
		pmClient, err := grpc.GetPrivateMessageClient()
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		reqUpdate := privatemessage.ReqPrivateMessageStatusUpdateDao{
			AccountId: userMeta.AccountId,
			Ids:       upm.MessageIds,
		}
		_, err = pmClient.PrivateMessageStatusUpdateDao(context.Background(), &reqUpdate)
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
	}
	handle.SuccessResp(c, "", nil)
}

// 私信内容查询
func PrivateMessageDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	messageId := c.Query("message_id")
	if messageId == "" {
		err := fmt.Errorf("message is isn't null")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	pmClient, err := grpc.GetPrivateMessageClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	reqDetail := privatemessage.ReqPrivateMessageDetailDao{
		AccountId: userMeta.AccountId,
		Id:        messageId,
	}
	var resp *privatemessage.RspPrivateMessageDetailDao
	resp, err = pmClient.PrivateMessageDetailDao(context.Background(), &reqDetail)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var pm = resp.PrivateMessage
	var userResp *user.RspUserDaoById
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{pm.SendId}})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if len(userResp.Users) != 1 {
		err = fmt.Errorf("query user error")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var detail = userResp.Users[0]
	var response = rsp.UserPrivateMessageDetailResp{
		UserPrivateMessage: rsp.PrivateMessagePbConvertToRsp(pm, detail),
		Content:            resp.Content,
	}
	handle.SuccessResp(c, "", response)
}

// 私信信息删除
func DeletePrivateMessageHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	ids := strings.Split(c.Query("ids"), ",")
	if len(ids) != 0 {
		pmClient, err := grpc.GetPrivateMessageClient()
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		req := &privatemessage.ReqPrivateMessageDeleteDao{
			AccountId: userMeta.AccountId,
			Ids:       ids,
		}
		_, err = pmClient.PrivateMessageDeleteDao(context.Background(), req)
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
	}
	handle.SuccessResp(c, "", nil)
}
