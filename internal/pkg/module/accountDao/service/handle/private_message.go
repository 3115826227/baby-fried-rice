package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"baby-fried-rice/internal/pkg/module/adminAccount/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendPrivateMessageHandle(c *gin.Context) {
	var pm requests.UserSendPrivateMessageReq
	if err := c.ShouldBind(&pm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	if err := db.SendPrivateMessage(pm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func PrivateMessagesHandle(c *gin.Context) {
	var pms requests.UserPrivateMessagesReq
	pms.UserId = c.Query("user_id")
	pageReq, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	pms.PageCommonReq = pageReq
	messages, err := query.GetUserPrivateMessages(pms)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var mp = make(map[string]tables.AccountUserDetail)
	for _, m := range messages {
		if _, exist := mp[m.SendId]; !exist {
			var detail tables.AccountUserDetail
			detail, err = query.GetUserDetail(m.SendId)
			if err != nil {
				log.Logger.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
			mp[m.SendId] = detail
		}
	}
	var resp = make([]rsp.UserPrivateMessagesResp, 0)
	for _, m := range messages {
		resp = append(resp, rsp.UserPrivateMessagesResp{
			MessageId:     m.MessageId,
			SendId:        m.SendId,
			SendName:      mp[m.SendId].Username,
			ReceiveId:     m.ReceiveId,
			MessageStatus: m.MessageStatus,
			ReceiveTime:   m.ReceiveTime.Format(config.TimeLayout),
		})
	}
	handle.SuccessResp(c, "", resp)
}

func UpdatePrivateMessageStatusHandle(c *gin.Context) {
	var upm requests.UpdatePrivateMessageStatusReq
	if err := c.ShouldBind(&upm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	if err := db.UpdatePrivateMessagesStatus(upm.ReceiveId, upm.MessageIds); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func PrivateMessageDetailHandle(c *gin.Context) {
	messageId := c.Query("message_id")
	var pmd tables.UserPrivateMessageContent
	if err := db.GetDB().GetObject(map[string]interface{}{"id": messageId}, &pmd); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", pmd)
}

func DeletePrivateMessageHandle(c *gin.Context) {
	var pm tables.UserPrivateMessage
	pm.ReceiveId = c.Query("receive_id")
	pm.MessageId = c.Query("message_id")
	if err := db.GetDB().DeleteObject(&pm); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
