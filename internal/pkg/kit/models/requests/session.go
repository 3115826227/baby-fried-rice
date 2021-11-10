package requests

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"fmt"
)

// 创建会话请求参数
type ReqAddSession struct {
	// 会话等级
	SessionLevel im.SessionLevel `json:"session_level"`
	// 会话类型
	SessionType im.SessionType `json:"session_type"`
	// 会话加入权限
	JoinPermissionType im.SessionJoinPermissionType `json:"join_permission_type"`
	// 会话名称
	Name string `json:"name"`
	// 加入会话成员id列表
	Joins []string `json:"joins"`
}

func (req *ReqAddSession) Validate() error {
	// 加入成员重复校验
	var joinMap = make(map[string]struct{})
	for _, user := range req.Joins {
		if _, exist := joinMap[user]; exist {
			return fmt.Errorf("join persion can't repeat")
		}
		joinMap[user] = struct{}{}
	}
	// 会话类型非双人必须要添加名称
	if req.SessionType != im.SessionType_DoubleSession && req.Name == "" {
		return fmt.Errorf("multi persion's session must have a name")
	}
	// 加入会话的成员人数不能超过会话等级限制的人数
	switch req.SessionLevel {
	case im.SessionLevel_SessionBaseLevel:
		if len(req.Joins) > 2 {
			return fmt.Errorf("base level session can't exceed two persion")
		}
	case im.SessionLevel_SessionNormalLevel:
		if len(req.Joins) > 2 {
			return fmt.Errorf("base level session can't exceed twity persion")
		}
	}
	return nil
}

// 会话信息更新请求参数
type ReqUpdateSession struct {
	SessionId          int64                        `json:"session_id"`
	JoinPermissionType im.SessionJoinPermissionType `json:"join_permission_type"`
	Name               string                       `json:"name"`
}

// 邀请加入会话请求参数
type ReqInviteJoinSession struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 邀请人id
	AccountId string `json:"account_id"`
}

// 从会话中移除请求参数
type ReqRemoveFromSession struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 被移除会话的成员id
	AccountId string `json:"account_id"`
}

// 用户操作添加请求参数
type ReqOperatorAdd struct {
	// 接收用户id
	Receive string `json:"receive"`
	// 操作类型
	OptType int64 `json:"opt_type"`
	// 操作内容
	Content string `json:"content"`
	// 是否需要确认
	NeedConfirm bool `json:"need_confirm"`
}

// 用户操作确认请求参数
type ReqOperatorConfirm struct {
	// 操作id
	OperatorId int64 `json:"operator_id"`
	// 确认结果
	Confirm bool `json:"confirm"`
}

// 用户操作读取状态更新请求参数
type ReqOperatorReadStatusUpdate struct {
	Operators []int64 `json:"operators"`
}

// 用户添加好友请求参数
type ReqAddFriend struct {
	// 好友id
	AccountId string `json:"account_id"`
	// 备注
	Remark string `json:"remark"`
}

// 用户黑名单更新请求参数
type ReqUpdateFriendBlackList struct {
	// 好友id
	Friend string `json:"friend"`
	// 加入黑名单/从黑名单返回好友列表
	BlackList bool `json:"black_list"`
}

// 用户备注更新请求参数
type ReqUpdateFriendRemark struct {
	// 好友id
	Friend string `json:"friend"`
	// 好友备注
	Remark string `json:"remark"`
}

// 用户管理信息更新请求参数
type ReqUserManageUpdate struct {
	AddFriendPermissionType int32 `json:"add_friend_permission_type"`
}

// 用户webrtc创建请求参数
type ReqCreateWebRTC struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 用户sdp
	Sdp string `json:"sdp"`
}

// 用户webrtc邀请加入参数
type ReqReturnWebRTC struct {
	// 回复
	Return bool `json:"return"`
	// 会话id
	SessionId int64 `json:"session_id"`
	// 用户sdp
	Sdp string `json:"sdp"`
}
