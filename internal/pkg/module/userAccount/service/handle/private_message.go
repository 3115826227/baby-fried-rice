package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"fmt"
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
	pmClient, err := grpc.GetPrivateMessageClient()
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
	_, err = pmClient.PrivateMessageAddDao(context.Background(), &reqAdd)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

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
	var idsMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		idsMap[u.Id] = rsp.User{
			AccountID: u.Id,
			Username:  u.Username,
		}
	}
	var list = make([]interface{}, 0)
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
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

func UpdatePrivateMessageStatusHandle(c *gin.Context) {
	var upm requests.UpdatePrivateMessageStatusReq
	if err := c.ShouldBind(&upm); err != nil {
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
	reqUpdate := privatemessage.ReqPrivateMessageStatusUpdateDao{
		AccountId: upm.AccountId,
		Ids:       upm.MessageIds,
	}
	_, err = pmClient.PrivateMessageStatusUpdateDao(context.Background(), &reqUpdate)
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
