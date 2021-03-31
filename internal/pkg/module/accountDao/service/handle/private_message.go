package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/req"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendPrivateMessageHandle(c *gin.Context) {

}

func PrivateMessagesHandle(c *gin.Context) {
	var pms req.UserPrivateMessagesReq
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
	var resp = make([]rsp.UserPrivateMessagesResp, 0)
	for _, m := range messages {
		resp = append(resp, rsp.UserPrivateMessagesResp{
			MessageId:     m.MessageId,
			SendId:        m.SendId,
			SendName:      "",
			ReceiveId:     m.ReceiveId,
			MessageStatus: m.MessageStatus,
		})
	}
}

func UpdatePrivateMessageStatusHandle(c *gin.Context) {

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
