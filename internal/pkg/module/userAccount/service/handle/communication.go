package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/errors"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func AddCommunicationHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.UserCommunicationAddReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *user.RspUserCommunicationAddDao
	var communicationReq = &user.ReqUserCommunicationAddDao{
		AccountId:         userMeta.AccountId,
		Title:             req.Title,
		CommunicationType: req.CommunicationType,
		Content:           req.Content,
	}
	if len(req.Images) != 0 {
		communicationReq.Images = strings.Join(req.Images, ",")
	}
	resp, err = userClient.UserCommunicationAddDao(context.Background(), communicationReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", resp.Id)
}

func CommunicationHandle(c *gin.Context) {
	pageReq, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	communicationTypeStr := c.Query("communication_type")
	var communicationType = 0
	if communicationTypeStr != "" {
		communicationType, err = strconv.Atoi(communicationTypeStr)
		if err != nil {
			log.Logger.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
			return
		}
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *user.RspUserCommunicationQueryDao
	resp, err = userClient.UserCommunicationQueryDao(context.Background(), &user.ReqUserCommunicationQueryDao{
		Origin:            userMeta.AccountId,
		CommunicationType: user.CommunicationType(communicationType),
		Page:              pageReq.Page,
		PageSize:          pageReq.PageSize,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, item := range resp.List {
		var communication = rsp.UserCommunicationResp{
			Id:                item.Id,
			Title:             item.Title,
			CommunicationType: item.CommunicationType,
			CreateTimestamp:   item.CreateTimestamp,
			UpdateTimestamp:   item.UpdateTimestamp,
			Reply:             item.Reply,
		}
		list = append(list, communication)
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

func CommunicationDetailHandle(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *user.RspUserCommunicationDetailQueryDao
	resp, err = userClient.UserCommunicationDetailQueryDao(context.Background(), &user.ReqUserCommunicationDetailQueryDao{
		Id:     int64(id),
		Origin: userMeta.AccountId,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(c, &user.ReqUserDaoById{Ids: []string{resp.Communication.Origin}})
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, err)
		return
	}
	if len(userResp.Users) != 1 {
		err = errors.NewCommonError(constant.CodeInternalError)
		log.Logger.Error(err.Error())
		handle.SystemErrorResponse(c)
		return
	}
	var origin = userResp.Users[0]
	var response = rsp.UserCommunicationDetailResp{
		UserCommunicationResp: rsp.UserCommunicationResp{
			Id:                resp.Communication.Id,
			Title:             resp.Communication.Title,
			CommunicationType: resp.Communication.CommunicationType,
			CreateTimestamp:   resp.Communication.CreateTimestamp,
			UpdateTimestamp:   resp.Communication.UpdateTimestamp,
		},
		Origin: rsp.User{
			AccountID:   origin.Id,
			Username:    origin.Username,
			HeadImgUrl:  origin.HeadImgUrl,
			IsOfficial:  origin.IsOfficial,
			PhoneVerify: origin.PhoneVerify,
		},
		Content:        resp.Content,
		Images:         resp.Images,
		ReplyContent:   resp.ReplyContent,
		ReplyTimestamp: resp.ReplyTimestamp,
	}
	handle.SuccessResp(c, "", response)
}

func DeleteCommunicationHandle(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	_, err = userClient.UserCommunicationDeleteDao(context.Background(), &user.ReqUserCommunicationDeleteDao{
		Id:     int64(id),
		Origin: userMeta.AccountId,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
