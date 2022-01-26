package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 空间操作添加
func OptAddHandle(c *gin.Context) {
	var req requests.ReqAddOpt
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	commentClient, err := grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpaceOpt = &comment.ReqOperatorAddDao{
		BizId:   req.BizId,
		BizType: req.BizType,
		HostId:  req.OperatorId,
		Origin:  userMeta.AccountId,
		OptType: req.OperatorType,
	}
	var resp *comment.RspOperatorAddDao
	resp, err = commentClient.OperatorAddDao(context.Background(), reqSpaceOpt)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	if !resp.Result {
		err = fmt.Errorf("operator add failed")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	} else {
		// 空间动态的操作需要自己处理
		if req.OperatorId == req.BizId {
			switch req.BizType {
			case comment.BizType_Space:
				if err = handleSpaceOpt(req); err != nil {
					c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
					return
				}
			}
		}
	}
	handle.SuccessResp(c, "", nil)
}

func handleSpaceOpt(req requests.ReqAddOpt) (err error) {
	var spaceClient space.DaoSpaceClient
	spaceClient, err = grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var spaceIncrReq = space.ReqSpaceIncrUpdateDao{
		Id: req.BizId,
	}
	switch req.OperatorType {
	case comment.OperatorType_Like:
		spaceIncrReq.LikeIncrement = 1
	case comment.OperatorType_CancelLike:
		spaceIncrReq.LikeIncrement = -1
	}
	if _, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &spaceIncrReq); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}
