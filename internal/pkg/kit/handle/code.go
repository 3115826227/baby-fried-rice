package handle

const (
	CodeInvalidParams     = 400
	CodeRequiredLogin     = 401
	CodeRequiredForbidden = 403
	CodeNotFound          = 404
	CodeInternalError     = 500
	CodeServiceNotFound   = 502

	CodeNeedApplyAddFriend = 20010

	CodeNeedOriginAuditSession = 20021
	CodeNeedInviteJoinSession  = 20022
)

const (
	CodeInvalidParamsMsg     = "参数错误"
	CodeRequiredLoginMsg     = "请登录"
	CodeRequiredForbiddenMsg = "权限不够"
	CodeNotFoundMsg          = "未找到服务"
	CodeInternalErrorMsg     = "服务器错误"
	CodeServiceNotFoundMsg   = "服务不存在"

	CodeNeedApplyAddFriendMsg = "对方已设置好友添加权限，请先申请添加好友"

	CodeNeedOriginAuditSessionMsg = "会话加入请求已发送，请耐心等待审核确认"
	CodeNeedInviteJoinSessionMsg  = "会话创建者已设置会话加入权限，请您联系会话创建者邀请您，才能加入该会话"
)
