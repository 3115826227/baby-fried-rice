package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func IteratorVersionHandle(c *gin.Context) {
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *user.RspIteratorVersionQueryDao
	resp, err = userClient.IteratorVersionQueryDao(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, v := range resp.List {
		list = append(list, rsp.IteratorVersionResp{
			Version:   v.Version,
			Content:   v.Content,
			Timestamp: v.Timestamp,
		})
	}
	handle.SuccessResp(c, "", list)
}
