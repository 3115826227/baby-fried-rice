package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/space/config"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
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
		VisitorType: int32(req.VisitorType),
	}
	_, err = client.SpaceAddDao(context.Background(), reqSpace)
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
		for _, id := range userResp.AccountIds {
			var notify = models.WSMessageNotify{
				WSMessageNotifyType: constant.SpaceMessageNotify,
				Receive:             id,
				WSMessage: models.WSMessage{
					Send: models.UserBaseInfo{
						AccountId:  detailResp.Detail.AccountId,
						Username:   detailResp.Detail.Username,
						HeadImgUrl: detailResp.Detail.HeadImgUrl,
					},
				},
				Timestamp: time.Now().Unix(),
			}
			if err = mq.Send(topic, notify.ToString()); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	}()
	handle.SuccessResp(c, "", nil)
}

func SpaceCommentConvert(comment *space.SpaceCommentDao, userClient user.DaoUserClient) (resp rsp.SpaceCommentResp, err error) {
	var comments = make([]rsp.SpaceCommentResp, 0)
	for _, reply := range comment.ReplyList {
		var cmt rsp.SpaceCommentResp
		if cmt, err = SpaceCommentConvert(reply, userClient); err != nil {
			return
		}
		comments = append(comments, cmt)
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{comment.Origin}})
	if err != nil {
		return
	}
	resp = rsp.SpaceCommentResp{
		ID:          comment.Id,
		SpaceId:     comment.SpaceId,
		Comment:     comment.Content,
		CommentType: comment.CommentType,
		CreateTime:  comment.CreateTime,
		Liked:       comment.Liked,
		Reply:       comments,
	}
	if len(userResp.Users) == 1 {
		resp.User = rsp.User{
			AccountID: userResp.Users[0].Id,
			Username:  userResp.Users[0].Username,
		}
	}
	return
}

/*
	空间查询接口
*/
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
	var resp *space.RspSpacesQueryDao
	resp, err = spaceClient.SpacesQueryDao(context.Background(),
		&space.ReqSpacesQueryDao{CommonSearchReq: searchReq})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	list := make([]rsp.SpaceResp, 0)
	for _, s := range resp.Spaces {
		var userResp *user.RspUserDaoById
		var ids = append(s.Other.Likes, s.Origin)
		var userClient user.DaoUserClient
		userClient, err = grpc.GetUserClient()
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		// 获取发布者用户名信息
		var originUser rsp.User
		for _, u := range userResp.Users {
			if s.Origin == u.Id {
				originUser = rsp.User{
					AccountID: u.Id,
					Username:  u.Username,
				}
				break
			}
		}
		var sp = rsp.SpaceResp{
			Id:          s.Id,
			Content:     s.Content,
			VisitorType: constant.SpaceVisitorType(s.VisitorType),
			Origin:      originUser,
			CreateTime:  s.CreateTime,
		}

		var likes = make([]rsp.User, 0)
		for _, u := range userResp.Users {
			likes = append(likes, rsp.User{
				AccountID: u.Id,
				Username:  u.Username,
			})
		}
		if s.Other != nil {
			sp.Other.Visited = s.Other.Visited
			sp.Other.Liked = s.Other.Liked
			sp.Other.Commented = s.Other.Commented
			sp.Other.Likes = likes
		}
		var comments = make([]rsp.SpaceCommentResp, 0)
		for _, cmt := range s.Other.Comments {
			var comment rsp.SpaceCommentResp
			if comment, err = SpaceCommentConvert(cmt, userClient); err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
			comments = append(comments, comment)
		}
		sp.Other.Comments = comments
		list = append(list, sp)
	}
	var response = rsp.SpacesResp{
		List:     list,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}
	handle.SuccessResp(c, "", response)
}

func SpaceDeleteHandle(c *gin.Context) {
	id := c.Query("id")
	spaceClient, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = spaceClient.SpaceDeleteDao(context.Background(), &space.ReqSpaceDeleteDao{Id: id})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func SpaceOptAddHandle(c *gin.Context) {
	var req requests.ReqAddSpaceOpt
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	spaceClient, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpaceOpt = &space.ReqSpaceOptAddDao{
		OperatorId:     req.OperatorId,
		SpaceId:        req.SpaceId,
		OperatorObject: req.OperatorObject,
		OperatorType:   req.OperatorType,
		Origin:         userMeta.AccountId,
	}
	_, err = spaceClient.SpaceOptAddDao(context.Background(), reqSpaceOpt)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func SpaceOptCancelHandle(c *gin.Context) {
	spaceId := c.Query("space_id")
	operatorId := c.Query("operator_id")
	userMeta := handle.GetUserMeta(c)
	spaceClient, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSpaceOpt = &space.ReqSpaceOptCancelDao{
		SpaceId:    spaceId,
		OperatorId: operatorId,
		Origin:     userMeta.AccountId,
	}
	_, err = spaceClient.SpaceOptCancelDao(context.Background(), reqSpaceOpt)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func SpaceCommentAddHandle(c *gin.Context) {
	var req requests.ReqAddSpaceComment
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	spaceClient, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userMeta := handle.GetUserMeta(c)
	var reqSpaceComment = &space.ReqSpaceCommentAddDao{
		SpaceId:     req.SpaceId,
		ParentId:    req.ParentId,
		Comment:     req.Comment,
		CommentType: req.CommentType,
		Origin:      userMeta.AccountId,
	}
	_, err = spaceClient.SpaceCommentAddDao(context.Background(), reqSpaceComment)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func SpaceCommentDeleteHandle(c *gin.Context) {
	id := c.Query("id")
	spaceId := c.Query("space_id")
	userMeta := handle.GetUserMeta(c)
	spaceClient, err := grpc.GetSpaceClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &space.ReqSpaceCommentDeleteDao{
		Id:      id,
		SpaceId: spaceId,
		Origin:  userMeta.AccountId,
	}
	_, err = spaceClient.SpaceCommentDeleteDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
