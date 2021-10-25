package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/space/config"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"time"
)

var (
	mq    interfaces.MQ
	topic string
)

func Init() {
	conf := config.GetConfig()
	topic = conf.MessageQueue.PublishTopics.WebsocketNotify
	mq = nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	if err := mq.NewProducer(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

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
		Images:      req.Images,
		VisitorType: req.VisitorType,
	}
	var resp *space.RspSpaceAddDao
	resp, err = client.SpaceAddDao(context.Background(), reqSpace)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		var userClient user.DaoUserClient
		userClient, err = grpc.GetUserClient()
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var detailResp *user.RspDaoUserDetail
		detailResp, err = userClient.UserDaoDetail(context.Background(), &user.ReqDaoUserDetail{AccountId: userMeta.AccountId})
		var userResp *user.RspUserDaoAll
		userResp, err = userClient.UserDaoAll(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var now = time.Now().Unix()
		for _, id := range userResp.AccountIds {
			var notify = models.WSMessageNotify{
				WSMessageNotifyType: constant.SpaceMessageNotify,
				Receive:             id,
				WSMessage: models.WSMessage{
					Space: &rsp.SpaceResp{
						Id:          resp.Id,
						Content:     req.Content,
						VisitorType: req.VisitorType,
						Origin: rsp.User{
							AccountID:  detailResp.Detail.AccountId,
							Username:   detailResp.Detail.Username,
							HeadImgUrl: detailResp.Detail.HeadImgUrl,
							IsOfficial: detailResp.Detail.IsOfficial,
						},
						CreateTime: now,
					},
					Send: models.UserBaseInfo{
						AccountId:  detailResp.Detail.AccountId,
						Username:   detailResp.Detail.Username,
						HeadImgUrl: detailResp.Detail.HeadImgUrl,
						IsOfficial: detailResp.Detail.IsOfficial,
					},
				},
				Timestamp: now,
			}
			if err = mq.Send(topic, notify.ToString()); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	}()
	handle.SuccessResp(c, "", resp.Id)
}

func findUsers(comments []*rsp.CommentResp) (ids []string) {
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

// 空间动态列表查询
func SpacesQueryHandle(c *gin.Context) {
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
		&space.ReqSpaceQueryDao{CommonSearchReq: searchReq, Origin: userMeta.AccountId})
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
		return
	}
	list := make([]rsp.SpaceResp, 0)
	for _, s := range resp.Spaces {
		var spe = rsp.SpaceResp{
			Id:           s.Id,
			Content:      s.Content,
			Images:       s.Images,
			VisitorType:  s.VisitorType,
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
		var commentResp *comment.RspCommentQueryDao
		if commentResp, err = commentClient.CommentQueryDao(context.Background(), &commentReq); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var comments = make([]*rsp.CommentResp, 0)
		for _, cmmt := range commentResp.List {
			comments = append(comments, rsp.CommentRpcConvertResponse(cmmt))
		}
		var commentIds = findUsers(comments)
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
	var commentClient comment.DaoCommentClient
	commentClient, err = grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	// 删除空间动态下的所有评论和操作
	_, err = commentClient.CommentClearDao(context.Background(), &comment.ReqCommentClearDao{BizId: id, BizType: comment.BizType_Space})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 空间操作添加
func SpaceOptAddHandle(c *gin.Context) {
	var req requests.ReqAddSpaceOpt
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
	var reqSpaceOpt = &comment.ReqOperatorAddDao{
		BizId:   req.SpaceId,
		BizType: comment.BizType_Space,
		HostId:  req.OperatorId,
		Origin:  userMeta.AccountId,
		OptType: req.OperatorType,
	}
	var resp *comment.RspOperatorAddDao
	resp, err = commentClient.OperatorAddDao(context.Background(), reqSpaceOpt)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if !resp.Result {
		err = fmt.Errorf("operator add failed")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	} else {
		// 空间动态的操作需要自己处理
		if req.OperatorId == req.SpaceId {
			var spaceClient space.DaoSpaceClient
			spaceClient, err = grpc.GetSpaceClient()
			if err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
			var spaceIncrReq = space.ReqSpaceIncrUpdateDao{
				Id: req.SpaceId,
			}
			switch req.OperatorType {
			case comment.OperatorType_Like:
				spaceIncrReq.LikeIncrement = 1
			case comment.OperatorType_CancelLike:
				spaceIncrReq.LikeIncrement = -1
			}
			if _, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &spaceIncrReq); err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
		}
	}
	handle.SuccessResp(c, "", nil)
}

// 空间动态评论添加
func SpaceCommentAddHandle(c *gin.Context) {
	var req requests.ReqAddSpaceComment
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
		BizId:    req.SpaceId,
		BizType:  comment.BizType_Space,
		ParentId: req.ParentId,
		Content:  req.Comment,
		Origin:   userMeta.AccountId,
	}
	var spaceIncrReq = space.ReqSpaceIncrUpdateDao{
		Id:               req.SpaceId,
		CommentIncrement: 1,
	}
	var spaceClient space.DaoSpaceClient
	spaceClient, err = grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if req.ParentId == "" {
		// 楼层评论，需要获取最新楼层并填入请求参数中，并更新空间信息
		var spaceResp *space.RspSpaceQueryDao
		spaceResp, err = spaceClient.SpaceQueryDao(context.Background(), &space.ReqSpaceQueryDao{SpaceId: req.SpaceId, Origin: userMeta.AccountId})
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		if len(spaceResp.Spaces) != 1 {
			err = fmt.Errorf("request space id is invalid")
			log.Logger.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
			return
		}
		reqComment.Floor = spaceResp.Spaces[0].FloorTotal + 1
		spaceIncrReq.FloorIncrement = 1
	}
	_, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &spaceIncrReq)
	if err != nil {
		log.Logger.Error(err.Error())
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

// 空间动态评论删除
func SpaceCommentDeleteHandle(c *gin.Context) {
	id := c.Query("id")
	bizId := c.Query("space_id")
	userMeta := handle.GetUserMeta(c)
	commentClient, err := grpc.GetCommentClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &comment.ReqCommentDeleteDao{
		Id:      id,
		BizType: comment.BizType_Space,
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
	// 更新空间动态评论总数值
	var spaceClient space.DaoSpaceClient
	spaceClient, err = grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = spaceClient.SpaceIncrUpdateDao(context.Background(), &space.ReqSpaceIncrUpdateDao{Id: req.BizId, CommentIncrement: -resp.Total})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
