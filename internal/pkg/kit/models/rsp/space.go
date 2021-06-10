package rsp

import "baby-fried-rice/internal/pkg/kit/constant"

type SpaceResp struct {
	Id          string                    `json:"id"`
	Content     string                    `json:"content"`
	VisitorType constant.SpaceVisitorType `json:"visitor_type"`
	Origin      string                    `json:"origin"`
	CreateTime  string                    `json:"create_time"`
	Other       SpaceOtherResp            `json:"other"`
}

type SpaceOtherResp struct {
	Visited   int64              `json:"visited"`
	Liked     int64              `json:"liked"`
	Commented int64              `json:"commented"`
	Likes     []User             `json:"likes"`
	Comments  []SpaceCommentResp `json:"comments"`
}

type SpacesResp struct {
	List     []SpaceResp `json:"list"`
	Page     int64       `json:"page"`
	PageSize int64       `json:"page_size"`
}

type SpaceCommentResp struct {
	ID          string             `json:"id"`
	SpaceId     string             `json:"space_id"`
	User        User               `json:"origin"`
	Comment     string             `json:"comment"`
	CommentType int32              `json:"comment_type"`
	CreateTime  string             `json:"create_time"`
	Reply       []SpaceCommentResp `json:"reply"`
}
