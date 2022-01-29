package rsp

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
)

type SpaceResp struct {
	Id                     string                 `json:"id"`
	Content                string                 `json:"content"`
	Images                 []string               `json:"images"`
	VisitorType            space.SpaceVisitorType `json:"visitor_type"`
	Origin                 *User                  `json:"origin,omitempty"`
	CreateTime             int64                  `json:"create_time"`
	VisitTotal             int64                  `json:"visit_total"`
	LikeTotal              int64                  `json:"like_total"`
	FloorTotal             int64                  `json:"floor_total"`
	OriginLiked            bool                   `json:"origin_liked"`
	CommentTotal           int64                  `json:"comment_total"`
	Forward                bool                   `json:"forward"`
	ForwardSpace           ForwardSpace           `json:"forward_space"`
	ForwardTotal           int64                  `json:"forward_total"`
	OpenComment            bool                   `json:"open_comment"`
	Comments               []*CommentResp         `json:"comments"`
	CurrentCommentPage     int64                  `json:"current_comment_page"`
	CurrentCommentPageSize int64                  `json:"current_comment_page_size"`
	Anonymity              bool                   `json:"anonymity"`
	OriginSpace            bool                   `json:"origin_space"`
}

type ForwardSpace struct {
	SpaceId     string                 `json:"space_id"`
	Content     string                 `json:"content"`
	Images      []string               `json:"images"`
	Origin      *User                  `json:"origin"`
	VisitorType space.SpaceVisitorType `json:"visitor_type"`
}

type SpacesResp struct {
	List     []SpaceResp `json:"list"`
	Page     int64       `json:"page"`
	PageSize int64       `json:"page_size"`
}

// 评论
type CommentResp struct {
	ID                   string       `json:"id"`
	User                 User         `json:"origin"`
	Content              string       `json:"content"`
	CreateTime           int64        `json:"create_time"`
	Floor                int64        `json:"floor"`
	LikeTotal            int64        `json:"like_total"`
	OriginLiked          bool         `json:"origin_liked"`
	OpenReply            bool         `json:"open_reply"`
	ReplyTotal           int64        `json:"reply_total"`
	Reply                []*ReplyResp `json:"reply"`
	CurrentReplyPage     int64        `json:"current_reply_page"`
	CurrentReplyPageSize int64        `json:"current_reply_page_size"`
}

func (response *CommentResp) SetUser(idMap map[string]User) {
	response.User = idMap[response.User.AccountID]
	for index, reply := range response.Reply {
		reply.User = idMap[reply.User.AccountID]
		reply.SetUser(idMap)
		response.Reply[index] = reply
	}
}

func (response *CommentResp) FindUserIds() map[string]User {
	var idMap = make(map[string]User)
	idMap[response.User.AccountID] = response.User
	for _, reply := range response.Reply {
		idMap[reply.User.AccountID] = reply.User
		var replyMap = reply.FindUserIds()
		for id, user := range replyMap {
			idMap[id] = user
		}
	}
	return idMap
}

func CommentRpcConvertResponse(comment *comment.CommentDao) *CommentResp {
	var replyResponse = make([]*ReplyResp, 0)
	for _, reply := range comment.Reply {
		replyResponse = append(replyResponse, ReplyRpcConvertResponse(reply))
	}
	return &CommentResp{
		ID: comment.Id,
		User: User{
			AccountID: comment.Origin,
		},
		Content:     comment.Content,
		CreateTime:  comment.CreateTimestamp,
		Floor:       comment.Floor,
		LikeTotal:   comment.LikeTotal,
		OriginLiked: comment.OriginLiked,
		ReplyTotal:  comment.ReplyTotal,
		Reply:       replyResponse,
	}
}

// 评论回复
type ReplyResp struct {
	ID          string       `json:"id"`
	ParentId    string       `json:"parent_id"`
	User        User         `json:"origin"`
	Content     string       `json:"content"`
	CreateTime  int64        `json:"create_time"`
	LikeTotal   int64        `json:"like_total"`
	OriginLiked bool         `json:"origin_liked"`
	OpenReply   bool         `json:"open_reply"`
	Reply       []*ReplyResp `json:"reply"`
}

func (response *ReplyResp) SetUser(idMap map[string]User) {
	response.User = idMap[response.User.AccountID]
	for index, reply := range response.Reply {
		reply.User = idMap[reply.User.AccountID]
		reply.SetUser(idMap)
		response.Reply[index] = reply
	}
}

func (response *ReplyResp) FindUserIds() map[string]User {
	var idMap = make(map[string]User)
	idMap[response.User.AccountID] = response.User
	for _, childReply := range response.Reply {
		idMap[childReply.User.AccountID] = childReply.User
		var childReplyMap = childReply.FindUserIds()
		for id, user := range childReplyMap {
			idMap[id] = user
		}
	}
	return idMap
}

func ReplyRpcConvertResponse(reply *comment.CommentReplyDao) *ReplyResp {
	var childResponses = make([]*ReplyResp, 0)
	for _, childReply := range reply.Reply {
		childResponses = append(childResponses, ReplyRpcConvertResponse(childReply))
	}
	return &ReplyResp{
		ID:       reply.Id,
		ParentId: reply.ParentId,
		User: User{
			AccountID: reply.Origin,
		},
		Content:     reply.Content,
		CreateTime:  reply.CreateTimestamp,
		LikeTotal:   reply.LikeTotal,
		OriginLiked: reply.OriginLiked,
		Reply:       childResponses,
	}
}

type AdminSpaceResp struct {
	Id          string                 `json:"id"`
	VisitorType space.SpaceVisitorType `json:"visitor_type"`
	AuditStatus int64                  `json:"audit_status"`
	User        User                   `json:"user"`
	CreateTime  int64                  `json:"create_time"`
	UpdateTime  int64                  `json:"update_time"`
}

type AdminSpaceDetailResp struct {
	AdminSpaceResp
	Content   string `json:"content"`
	Visited   int64  `json:"visited"`
	Liked     int64  `json:"liked"`
	Commented int64  `json:"commented"`
}
