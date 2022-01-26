package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/module/im/grpc"
	"baby-fried-rice/internal/pkg/module/im/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ImgCollectAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqAddUserImgCollect
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInvalidParams)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	_, err = imClient.UserImgCollectAddDao(c, &im.ReqUserImgCollectAddDao{
		Img:       req.Img,
		AccountId: userMeta.AccountId,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func ImgCollectQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var imClient im.DaoImClient
	imClient, err = grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	var resp *im.RspUserImgCollectQueryDao
	resp, err = imClient.UserImgCollectQueryDao(c, &im.ReqUserImgCollectQueryDao{
		AccountId: userMeta.AccountId,
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	var list = make([]interface{}, 0)
	for _, img := range resp.List {
		list = append(list, img)
	}
	handle.SuccessListResp(c, "", list, resp.Total, reqPage.Page, reqPage.PageSize)
}

func ImgCollectDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	img := c.Query("img")
	if img == "" {
		handle.FailedResp(c, constant.CodeInvalidParams)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	_, err = imClient.UserImgCollectDeleteDao(c, &im.ReqUserImgCollectDeleteDao{
		AccountId: userMeta.AccountId,
		Img:       img,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	handle.SuccessResp(c, "", nil)
}
