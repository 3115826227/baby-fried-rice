package requests

// 创建会话请求参数
type ReqAddSession struct {
	// 会话等级
	SessionLevel int64 `json:"session_level"`
	// 会话类型
	SessionType int32 `json:"session_type"`
	// 会话加入权限
	JoinPermissionType int32 `json:"join_permission_type"`
	// 会话名称
	Name string `json:"name"`
	// 加入会话成员id列表
	Joins []string `json:"joins"`
}

// 会话信息更新请求参数
type ReqUpdateSession struct {
	SessionId          int64  `json:"session_id"`
	JoinPermissionType int32  `json:"join_permission_type"`
	Name               string `json:"name"`
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
