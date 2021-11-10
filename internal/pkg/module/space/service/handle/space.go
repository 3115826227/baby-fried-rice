package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 空间动态添加
func SpaceAddHandle(c *gin.Context) {
	var req requests.ReqAddSpace
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	client, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpace = &space.ReqSpaceAddDao{
		Origin:      userMeta.AccountId,
		Content:     req.Content,
		VisitorType: req.VisitorType,
	}
	if len(req.Images) != 0 {
		reqSpace.Images = req.Images
	}
	var resp *space.RspSpaceAddDao
	resp, err = client.SpaceAddDao(context.Background(), reqSpace)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		var spe = rsp.SpaceResp{
			Id:          resp.Id,
			Content:     req.Content,
			Images:      req.Images,
			VisitorType: req.VisitorType,
		}
		sendAddSpaceNotify(spe, userMeta.AccountId)
	}()
	handle.SuccessResp(c, "", resp.Id)
}

func findCommentUsers(comments []*rsp.CommentResp) (ids []string) {
	var idsMap = make(map[string]rsp.User)
	for _, c := range comments {
		commentIdsMap := c.FindUserIds()
		for id, u := range commentIdsMap {
			idsMap[id] = u
		}
	}
	for id := range idsMap {
		ids = append(ids, id)
	}
	return
}

func findReplyUsers(replies []*rsp.ReplyResp) (ids []string) {
	var idsMap = make(map[string]rsp.User)
	for _, reply := range replies {
		replyIdsMap := reply.FindUserIds()
		for id, u := range replyIdsMap {
			idsMap[id] = u
		}
	}
	for id := range idsMap {
		ids = append(ids, id)
	}
	return
}

// 空间动态列表查询
func SpacesQueryHandle(c *gin.Context) {
	visitorType, err := strconv.Atoi(c.Query("visitor_type"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req, err := handle.PageHandle(c)
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var spaceClient space.DaoSpaceClient
	spaceClient, err = grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	searchReq := &common.CommonSearchRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	userMeta := handle.GetUserMeta(c)
	var resp *space.RspSpaceQueryDao
	resp, err = spaceClient.SpaceQueryDao(context.Background(),
		&space.ReqSpaceQueryDao{CommonSearchReq: searchReq, Origin: userMeta.AccountId, VisitorType: space.SpaceVisitorType(visitorType)})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var idsMap = make(map[string]rsp.User)
	for _, s := range resp.Spaces {
		idsMap[s.Origin] = rsp.User{}
	}
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var ids []string
	for id := range idsMap {
		ids = append(ids, id)
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	for _, u := range userResp.Users {
		idsMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var commentClient comment.DaoCommentClient
	commentClient, err = grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	list := make([]rsp.SpaceResp, 0)
	for _, s := range resp.Spaces {
		var spe = rsp.SpaceResp{
			Id:           s.Id,
			Content:      s.Content,
			VisitorType:  s.VisitorType,
			Images:       s.Images,
			Origin:       idsMap[s.Origin],
			CreateTime:   s.CreateTime,
			VisitTotal:   s.VisitTotal,
			LikeTotal:    s.LikeTotal,
			FloorTotal:   s.FloorTotal,
			OriginLiked:  s.OriginLiked,
			CommentTotal: s.CommentTotal,
		}
		var commentReq = comment.ReqCommentQueryDao{
			BizId:    s.Id,
			BizType:  comment.BizType_Space,
			Origin:   userMeta.AccountId,
			Page:     1,
			PageSize: 5,
		}
		var comments []*rsp.CommentResp
		comments, _, err = commentQueryHandle(commentReq)
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		spe.Comments = comments
		list = append(list, spe)
	}
	go func() {
		// 更新空间浏览数
		for _, s := range resp.Spaces {
			var visitReq = comment.ReqVisitAddDao{
				BizId:     s.Id,
				BizType:   comment.BizType_Space,
				AccountId: userMeta.AccountId,
			}
			var visitResp *comment.RspVisitAddDao
			visitResp, err = commentClient.VisitAddDao(context.Background(), &visitReq)
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			if visitResp.Result {
				_, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &space.ReqSpaceIncrUpdateDao{Id: s.Id, VisitIncrement: 1})
				if err != nil {
					log.Logger.Error(err.Error())
					continue
				}
			}
		}
	}()
	var response = rsp.SpacesResp{
		List:     list,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}
	handle.SuccessResp(c, "", response)
}

// 空间动态删除
func SpaceDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	id := c.Query("id")
	spaceClient, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	// 删除空间动态
	_, err = spaceClient.SpaceDeleteDao(context.Background(), &space.ReqSpaceDeleteDao{Id: id, Origin: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if err = handleClearComment(id, comment.BizType_Space); err != nil {
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
