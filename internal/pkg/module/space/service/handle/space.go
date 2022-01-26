package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
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
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// 空间动态添加
func SpaceAddHandle(c *gin.Context) {
	var req requests.ReqAddSpace
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	// 请求参数校验
	if err := req.Validate(); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, err.Code())
		return
	}
	client, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpace = &space.ReqSpaceAddDao{
		Origin:      userMeta.AccountId,
		Content:     req.Content,
		VisitorType: req.VisitorType,
		Anonymity:   req.Anonymity,
	}
	if len(req.Images) != 0 {
		reqSpace.Images = req.Images
	}
	var resp *space.RspSpaceAddDao
	resp, err = client.SpaceAddDao(context.Background(), reqSpace)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, constant.CodeInternalError)
		return
	}
	go func() {
		var spe = rsp.SpaceResp{
			Id:          resp.Id,
			Content:     req.Content,
			Images:      req.Images,
			VisitorType: req.VisitorType,
			Anonymity:   req.Anonymity,
		}
		sendAddSpaceNotify(spe, userMeta.AccountId)
	}()
	handle.SuccessResp(c, "", resp.Id)
}

// 空间动态转发
func SpaceForwardHandle(c *gin.Context) {
	var req requests.ReqForwardSpace
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	client, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpace = &space.ReqSpaceForwardDao{
		Origin:      userMeta.AccountId,
		SpaceId:     req.OriginSpaceId,
		Content:     req.Content,
		VisitorType: req.VisitorType,
	}
	var resp *space.RspSpaceForwardDao
	resp, err = client.SpaceForwardDao(context.Background(), reqSpace)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.SystemErrorResponse(c)
		return
	}
	go func() {
		var spaceBaseResp *space.RspSpaceBaseQueryDao
		spaceBaseResp, err = client.SpaceBaseQueryDao(c, &space.ReqSpaceBaseQueryDao{Ids: []string{req.OriginSpaceId}})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var fs rsp.ForwardSpace
		if len(spaceBaseResp.List) == 1 {
			var forwardSpaceResp = spaceBaseResp.List[0]
			var users map[string]*rsp.User
			users, err = getUserByIds([]string{forwardSpaceResp.Origin})
			if err != nil {
				log.Logger.Error(err.Error())
				return
			}
			fs = rsp.ForwardSpace{
				SpaceId: forwardSpaceResp.Id,
				Content: forwardSpaceResp.Content,
				Images:  forwardSpaceResp.Images,
				Origin:  users[forwardSpaceResp.Origin],
			}
		}
		var spe = rsp.SpaceResp{
			Id:           resp.Id,
			Content:      req.Content,
			VisitorType:  req.VisitorType,
			Forward:      true,
			ForwardSpace: fs,
		}
		sendAddSpaceNotify(spe, userMeta.AccountId)
	}()
	handle.SuccessResp(c, "", resp.Id)
}

func getUserByIds(ids []string) (users map[string]*rsp.User, err error) {
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		err = errors.Wrap(err, "failed to get user client")
		return
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		err = errors.Wrap(err, "failed to get user by id")
		return
	}
	users = make(map[string]*rsp.User)
	for _, u := range userResp.Users {
		users[u.Id] = &rsp.User{
			AccountID:   u.Id,
			Username:    u.Username,
			HeadImgUrl:  u.HeadImgUrl,
			IsOfficial:  u.IsOfficial,
			PhoneVerify: u.PhoneVerify,
		}
	}
	return
}

func getForwardSpace(ids []string, spaceClient space.DaoSpaceClient) (forwardSpaces map[string]rsp.ForwardSpace, err error) {
	var baseSpaceResp *space.RspSpaceBaseQueryDao
	baseSpaceResp, err = spaceClient.SpaceBaseQueryDao(context.Background(), &space.ReqSpaceBaseQueryDao{Ids: ids})
	if err != nil {
		err = errors.Wrap(err, "failed to get space base query")
		return
	}
	forwardSpaces = make(map[string]rsp.ForwardSpace)
	var userIds = make([]string, 0)
	for _, s := range baseSpaceResp.List {
		userIds = append(userIds, s.Id)
	}
	var users map[string]*rsp.User
	users, err = getUserByIds(userIds)
	if err != nil {
		err = errors.Wrap(err, "failed to get user by ids")
		return
	}
	for _, s := range baseSpaceResp.List {
		forwardSpaces[s.Id] = rsp.ForwardSpace{
			SpaceId:     s.Id,
			Content:     s.Content,
			Images:      s.Images,
			Origin:      users[s.Origin],
			VisitorType: s.VisitorType,
		}
	}
	return
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
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var req requests.PageCommonReq
	req, err = handle.PageHandle(c)
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var spaceClient space.DaoSpaceClient
	spaceClient, err = grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		handle.SystemErrorResponse(c)
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
		handle.SystemErrorResponse(c)
		return
	}
	var idsMap = make(map[string]*rsp.User)
	for _, s := range resp.Spaces {
		if s.Anonymity {
			// 过滤掉匿名动态的用户
			continue
		}
		idsMap[s.Origin] = &rsp.User{}
	}
	var ids []string
	for id := range idsMap {
		ids = append(ids, id)
	}
	idsMap, err = getUserByIds(ids)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.SystemErrorResponse(c)
		return
	}
	var commentClient comment.DaoCommentClient
	commentClient, err = grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		handle.SystemErrorResponse(c)
		return
	}
	list := make([]rsp.SpaceResp, 0)
	var forwardSpaceIds = make([]string, 0)
	for _, s := range resp.Spaces {
		if s.Forward {
			forwardSpaceIds = append(forwardSpaceIds, s.OriginSpaceId)
		}
	}
	var forwardSpaces map[string]rsp.ForwardSpace
	forwardSpaces, err = getForwardSpace(forwardSpaceIds, spaceClient)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.SystemErrorResponse(c)
		return
	}
	for _, s := range resp.Spaces {
		var spe = rsp.SpaceResp{
			Id:           s.Id,
			Content:      s.Content,
			VisitorType:  s.VisitorType,
			Images:       s.Images,
			CreateTime:   s.CreateTime,
			VisitTotal:   s.VisitTotal,
			LikeTotal:    s.LikeTotal,
			FloorTotal:   s.FloorTotal,
			OriginLiked:  s.OriginLiked,
			CommentTotal: s.CommentTotal,
			Forward:      s.Forward,
			ForwardTotal: s.ForwardTotal,
			Anonymity:    s.Anonymity,
		}
		if !s.Anonymity {
			spe.Origin = idsMap[s.Origin]
		} else if s.Origin == userMeta.AccountId {
			spe.OriginSpace = true
		}
		if s.Forward {
			spe.ForwardSpace = forwardSpaces[s.OriginSpaceId]
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
			handle.SystemErrorResponse(c)
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
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	// 删除空间动态
	_, err = spaceClient.SpaceDeleteDao(context.Background(), &space.ReqSpaceDeleteDao{Id: id, Origin: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	if err = handleClearComment(id, comment.BizType_Space); err != nil {
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
