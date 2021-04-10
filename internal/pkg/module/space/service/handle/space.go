package handle

import (
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/space"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/space/config"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddSpaceHandle(c *gin.Context) {
	var err error
	var req requests.ReqAddSpace
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}

	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.SpaceDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpace = &space.ReqSpaceAddDao{
		Origin:      userMeta.UserId,
		Content:     req.Content,
		VisitorType: int32(req.VisitorType),
	}
	resp, err := space.NewDaoSpaceClient(client.GetRpcClient()).
		SpaceAddDao(context.Background(), reqSpace)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if resp.Code != handle.SuccessCode {
		log.Logger.Error(resp.Message)
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func QuerySpacesHandle(c *gin.Context) {

}

func DeleteSpaceHandle(c *gin.Context) {

}
