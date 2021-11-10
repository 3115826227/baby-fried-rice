package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 评论添加
func CommentAddHandle(c *gin.Context) {
	var req requests.ReqAddComment
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	commentClient, err := grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqComment = &comment.ReqCommentAddDao{
		BizId:    req.BizId,
		BizType:  comment.BizType_Space,
		ParentId: req.ParentId,
		Content:  req.Comment,
		Origin:   userMeta.AccountId,
	}
	if reqComment.Floor, err = handleBizCommentAdd(req, userMeta.AccountId); err != nil {
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *comment.RspCommentAddDao
	resp, err = commentClient.CommentAddDao(context.Background(), reqComment)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", resp.Id)
}

// 空间更多评论查询
func CommentQueryHandle(c *gin.Context) {
	bizId := c.Query("biz_id")
	bizType, err := strconv.Atoi(c.Query("biz_type"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var pageReq requests.PageCommonReq
	if pageReq, err = handle.PageHandle(c); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var req = comment.ReqCommentQueryDao{
		BizId:    bizId,
		BizType:  comment.BizType(bizType),
		Origin:   userMeta.AccountId,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}
	var list []interface{}
	comments, total, err := commentQueryHandle(req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	for _, cmmt := range comments {
		list = append(list, cmmt)
	}
	handle.SuccessListResp(c, "", list, total, pageReq.Page, pageReq.PageSize)
}

// 空间评论更多回复查询
func CommentReplyQueryHandle(c *gin.Context) {
	floor, err := strconv.Atoi(c.Query("floor"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	bizId := c.Query("biz_id")
	var bizType int
	if bizType, err = strconv.Atoi(c.Query("biz_type")); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var pageReq requests.PageCommonReq
	if pageReq, err = handle.PageHandle(c); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var req = comment.ReqCommentReplyQueryDao{
		BizId:    bizId,
		BizType:  comment.BizType(bizType),
		Floor:    int64(floor),
		Origin:   userMeta.AccountId,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	}
	var list []interface{}
	replies, total, err := commentReplyQueryHandle(req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	for _, reply := range replies {
		list = append(list, reply)
	}
	handle.SuccessListResp(c, "", list, total, pageReq.Page, pageReq.PageSize)
}

// 空间动态评论删除
func CommentDeleteHandle(c *gin.Context) {
	id := c.Query("id")
	bizId := c.Query("space_id")
	bizTypeStr := c.Query("biz_type")
	bizType, err := strconv.Atoi(bizTypeStr)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var commentClient comment.DaoCommentClient
	if commentClient, err = grpc.GetCommentClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &comment.ReqCommentDeleteDao{
		Id:      id,
		BizType: comment.BizType(bizType),
		BizId:   bizId,
		Origin:  userMeta.AccountId,
	}
	var resp *comment.RspCommentDeleteDao
	resp, err = commentClient.CommentDeleteDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if err = handleBizCommentTotal(bizId, comment.BizType(bizType), -resp.Total); err != nil {
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func commentQueryHandle(req comment.ReqCommentQueryDao) (comments []*rsp.CommentResp, total int64, err error) {
	var commentClient comment.DaoCommentClient
	commentClient, err = grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var commentResp *comment.RspCommentQueryDao
	if commentResp, err = commentClient.CommentQueryDao(context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	total = commentResp.Total
	for _, cmmt := range commentResp.List {
		comments = append(comments, rsp.CommentRpcConvertResponse(cmmt))
	}
	var commentIds = findCommentUsers(comments)
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.
		UserDaoById(context.Background(),
			&user.ReqUserDaoById{Ids: commentIds}); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var commentIdsMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		commentIdsMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	for index, cmmt := range comments {
		cmmt.SetUser(commentIdsMap)
		comments[index] = cmmt
	}
	return
}

func commentReplyQueryHandle(req comment.ReqCommentReplyQueryDao) (replies []*rsp.ReplyResp, total int64, err error) {
	var commentClient comment.DaoCommentClient
	commentClient, err = grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var replyResp *comment.RspCommentReplyQueryDao
	if replyResp, err = commentClient.CommentReplyQueryDao(context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	total = replyResp.Total
	for _, reply := range replyResp.List {
		replies = append(replies, rsp.ReplyRpcConvertResponse(reply))
	}
	var commentIds = findReplyUsers(replies)
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.
		UserDaoById(context.Background(),
			&user.ReqUserDaoById{Ids: commentIds}); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var commentIdsMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		commentIdsMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	for index, reply := range replies {
		reply.SetUser(commentIdsMap)
		replies[index] = reply
	}
	return
}

func handleBizCommentAdd(req requests.ReqAddComment, accountId string) (floor int64, err error) {
	switch req.BizType {
	case comment.BizType_Space:
		var spaceClient space.DaoSpaceClient
		spaceClient, err = grpc.GetSpaceClient()
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var spaceIncrReq = space.ReqSpaceIncrUpdateDao{
			Id:               req.BizId,
			CommentIncrement: 1,
		}
		if req.ParentId == "" {
			// 楼层评论，需要获取最新楼层并填入请求参数中，并更新空间信息
			var spaceResp *space.RspSpaceQueryDao
			spaceResp, err = spaceClient.SpaceQueryDao(context.Background(), &space.ReqSpaceQueryDao{SpaceId: req.BizId, Origin: accountId})
			if err != nil {
				log.Logger.Error(err.Error())
				return
			}
			if len(spaceResp.Spaces) != 1 {
				err = fmt.Errorf("request space id is invalid")
				log.Logger.Error(err.Error())
				return
			}
			floor = spaceResp.Spaces[0].FloorTotal + 1
			spaceIncrReq.FloorIncrement = 1
		}
		_, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &spaceIncrReq)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	return
}

func handleBizCommentTotal(bizId string, bizType comment.BizType, total int64) (err error) {
	switch bizType {
	case comment.BizType_Space:
		// 更新空间动态评论总数值
		var spaceClient space.DaoSpaceClient
		spaceClient, err = grpc.GetSpaceClient()
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		_, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &space.ReqSpaceIncrUpdateDao{Id: bizId, CommentIncrement: total})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	return
}

func handleClearComment(bizId string, bizType comment.BizType) (err error) {
	var commentClient comment.DaoCommentClient
	commentClient, err = grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 删除该业务类型下的所有评论和操作
	_, err = commentClient.CommentClearDao(context.Background(), &comment.ReqCommentClearDao{BizId: bizId, BizType: bizType})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}
