package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendPrivateMessageHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var pm requests.UserSendPrivateMessageReq
	pm.SendId = userMeta.AccountId
	if err := c.ShouldBind(&pm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	reqAdd := privatemessage.ReqPrivateMessageAddDao{
		SendId:          pm.SendId,
		ReceiveId:       pm.ReceiveId,
		MessageType:     pm.MessageType,
		MessageSendType: int32(pm.MessageSendType),
		Title:           pm.MessageTitle,
		Content:         pm.MessageContent,
	}
	_, err = privatemessage.NewDaoPrivateMessageClient(client.GetRpcClient()).
		PrivateMessageAddDao(context.Background(), &reqAdd)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func PrivateMessagesHandle(c *gin.Context) {
	sendId := c.Query("send_id")
	pageReq, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	reqQuery := privatemessage.ReqPrivateMessageQueryDao{
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
		SendId:   sendId,
	}
	var resp *privatemessage.RspPrivateMessageQueryDao
	resp, err = privatemessage.NewDaoPrivateMessageClient(client.GetRpcClient()).
		PrivateMessageQueryDao(context.Background(), &reqQuery)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]rsp.UserPrivateMessage, 0)
	var ids = make([]string, 0)
	for _, pm := range resp.List {
		ids = append(ids, pm.SendId)
	}
	var userResp *user.RspUserDaoById
	userResp, err = user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		return
	}
	var idsMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		idsMap[u.Id] = rsp.User{
			AccountID: u.Id,
			Username:  u.Username,
		}
	}
	for _, pm := range resp.List {
		var pmsg = rsp.UserPrivateMessage{
			MessageId:     pm.Id,
			Send:          idsMap[pm.SendId],
			ReceiveId:     pm.ReceiveId,
			MessageStatus: pm.Status,
			ReceiveTime:   pm.CreateTime,
			Title:         pm.Title,
		}
		list = append(list, pmsg)
	}
	var response = rsp.UserPrivateMessagesResp{
		List:     list,
		Page:     resp.Page,
		PageSize: resp.PageSize,
		Total:    resp.Total,
	}
	handle.SuccessResp(c, "", response)
}

func UpdatePrivateMessageStatusHandle(c *gin.Context) {
	var upm requests.UpdatePrivateMessageStatusReq
	if err := c.ShouldBind(&upm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	reqUpdate := privatemessage.ReqPrivateMessageStatusUpdateDao{
		AccountId: upm.AccountId,
		Ids:       upm.MessageIds,
	}
	_, err = privatemessage.NewDaoPrivateMessageClient(client.GetRpcClient()).
		PrivateMessageStatusUpdateDao(context.Background(), &reqUpdate)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func PrivateMessageDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	messageId := c.Query("message_id")
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
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
	resp, err = privatemessage.NewDaoPrivateMessageClient(client.GetRpcClient()).
		PrivateMessageDetailDao(context.Background(), &reqDetail)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var pm = resp.PrivateMessage
	var userResp *user.RspUserDaoById
	userResp, err = user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{pm.SendId}})
	if err != nil {
		return
	}
	if len(userResp.Users) != 1 {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var response = rsp.UserPrivateMessageDetailResp{
		UserPrivateMessage: rsp.UserPrivateMessage{
			MessageId: pm.Id,
			Send: rsp.User{
				AccountID: userResp.Users[0].Id,
				Username:  userResp.Users[0].Username,
			},
			ReceiveId:     pm.ReceiveId,
			MessageStatus: pm.Status,
			ReceiveTime:   pm.CreateTime,
			Title:         pm.Title,
		},
		Content: resp.Content,
	}
	handle.SuccessResp(c, "", response)
}

func DeletePrivateMessageHandle(c *gin.Context) {
}
