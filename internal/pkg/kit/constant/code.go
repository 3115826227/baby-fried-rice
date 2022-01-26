package constant

import "github.com/gin-gonic/gin"

type Code uint32

const (
	CodeInvalidParams     Code = 400
	CodeRequiredLogin          = 401
	CodeRequiredForbidden      = 403
	CodeNotFound               = 404
	CodeInternalError     Code = 500
	CodeServiceNotFound        = 502
	CodeUnVerifyForbidden      = 600

	CodeNeedApplyAddFriend = 20010

	CodeNeedOriginAuditSession = 20021
	CodeNeedInviteJoinSession  = 20022

	CodeSessionOriginPermission = 20031

	ErrCodeLoginFailed    = 99
	CodeLoginNameEmpty    = 30011
	CodeLoginNameExist    = 30012
	CodePasswordEmpty     = 30013
	CodeLoginNameNotExist = 30014
	CodePasswordInvalid   = 30015
	CodeUsernameEmpty     = 30016
	CodeUserFreeze        = 30017
	CodeUserCancel        = 30018

	CodePhoneInvalid     = 20041
	CodePhoneEmpty       = 20042
	CodePhoneVerifyExist = 20043

	CodePhoneVerifyCodeTooBusy = 20011
	CodePhoneVerifyCodeError   = 20012
	CodePhoneVerifyCodeInvalid = 20013
	CodePhoneVerifyCodeExpire  = 20014
	CodePhoneVerifyCodeEmpty   = 20015

	CodeSelfVideoConflictError = 20101
	CodeUserVideoConflictError = 20102
	CodeUserOfflineError       = 20103

	CodeSpaceContentEmptyError       = 20501
	CodeSpaceVisitorTypeInvalidError = 20511
)

const (
	CodeInvalidParamsMsg     = "参数错误"
	CodeRequiredLoginMsg     = "请登录"
	CodeRequiredForbiddenMsg = "权限不够"
	CodeNotFoundMsg          = "未找到服务"
	CodeInternalErrorMsg     = "服务器错误"
	CodeServiceNotFoundMsg   = "服务不存在"
	CodeUnVerifyForbiddenMsg = "未认证用户无法访问"

	CodeNeedApplyAddFriendMsg = "对方已设置好友添加权限，请先申请添加好友"

	CodeNeedOriginAuditSessionMsg = "会话加入请求已发送，请耐心等待审核确认"
	CodeNeedInviteJoinSessionMsg  = "会话创建者已设置会话加入权限，请您联系会话创建者邀请您，才能加入该会话"

	CodeSessionOriginPermissionMsg = "只有会话创建才有该权限"

	// 用户相关
	ErrCodeLoginFailedMsg    = "用户名或密码错误"
	CodeLoginNameEmptyMsg    = "登录账号不能为空"
	CodeLoginNameExistMsg    = "登录账号已存在"
	CodePasswordEmptyMsg     = "登录密码不能为空"
	CodeLoginNameNotExistMsg = "登录账号不存在"
	CodePasswordInvalidMsg   = "无效的登录密码"
	CodeUsernameEmptyMsg     = "用户名不能为空"
	CodeUserFreezeMsg        = "用户账号被冻结"
	CodeUserCancelMsg        = "用户账号已经注销"

	// 手机号相关
	CodePhoneInvalidMsg     = "无效的手机号码"
	CodePhoneEmptyMsg       = "手机号不能为空"
	CodePhoneVerifyExistMsg = "手机号已经被验证"

	// 验证码相关
	CodePhoneVerifyCodeTooBusyMsg = "验证码发送太频繁，请稍后再试"
	CodePhoneVerifyCodeErrorMsg   = "验证码申请失败，请重试"
	CodePhoneVerifyCodeInvalidMsg = "验证码无效，请填写正确的验证码"
	CodePhoneVerifyCodeExpireMsg  = "找不到验证的手机号，请填写正确的手机号，重新获取短信验证码"
	CodePhoneVerifyCodeEmptyMsg   = "验证码不能为空，请填写正确有效的验证码"

	CodeSelfVideoConflictErrorMsg = "您有未结束的通话，请先结束当前会话"
	CodeUserVideoConflictErrorMsg = "对方正在通话，请稍后再试"
	CodeUserOfflineErrorMsg       = "对方未上线，请稍后再试"

	CodeSpaceContentEmptyErrorMsg       = "动态内容不能为空"
	CodeSpaceVisitorTypeInvalidErrorMsg = "动态访问类型参数非法"
)

const (
	InternalCodePhoneEmptyMsg   = "phone is empty"
	InternalCodePhoneInvalidMsg = "phone is invalid"
)

var ErrCodeM = map[Code]string{
	ErrCodeLoginFailed:    ErrCodeLoginFailedMsg,
	CodeInvalidParams:     CodeInvalidParamsMsg,
	CodeInternalError:     CodeInternalErrorMsg,
	CodeRequiredForbidden: CodeRequiredForbiddenMsg,

	CodeLoginNameEmpty:    CodeLoginNameEmptyMsg,
	CodeLoginNameExist:    CodeLoginNameExistMsg,
	CodePasswordEmpty:     CodePasswordEmptyMsg,
	CodeLoginNameNotExist: CodeLoginNameNotExistMsg,
	CodePasswordInvalid:   CodePasswordInvalidMsg,
	CodeUsernameEmpty:     CodeUsernameEmptyMsg,
	CodeUserFreeze:        CodeUserFreezeMsg,
	CodeUserCancel:        CodeUserCancelMsg,

	CodePhoneInvalid:           CodePhoneInvalidMsg,
	CodePhoneEmpty:             CodePhoneEmptyMsg,
	CodePhoneVerifyExist:       CodePhoneVerifyExistMsg,
	CodePhoneVerifyCodeTooBusy: CodePhoneVerifyCodeTooBusyMsg,
	CodePhoneVerifyCodeError:   CodePhoneVerifyCodeErrorMsg,
	CodePhoneVerifyCodeInvalid: CodePhoneVerifyCodeInvalidMsg,
	CodePhoneVerifyCodeExpire:  CodePhoneVerifyCodeExpireMsg,
	CodePhoneVerifyCodeEmpty:   CodePhoneVerifyCodeEmptyMsg,
	CodeSelfVideoConflictError: CodeSelfVideoConflictErrorMsg,
	CodeUserVideoConflictError: CodeUserVideoConflictErrorMsg,
	CodeUserOfflineError:       CodeUserOfflineErrorMsg,

	CodeSpaceContentEmptyError:       CodeSpaceContentEmptyErrorMsg,
	CodeSpaceVisitorTypeInvalidError: CodeSpaceVisitorTypeInvalidErrorMsg,
}

var ParamErrResponse = gin.H{
	"code": CodeInvalidParams,
	"data": make(map[string]interface{}),
}

var SysErrResponse = gin.H{
	"code": CodeInternalError,
	"data": make(map[string]interface{}),
}
