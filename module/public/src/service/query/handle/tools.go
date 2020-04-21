package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrCodeLoginFailed              = 99
	ErrCodeLoginPlease              = 100
	ErrCodeUserBanned               = 101
	ErrCodeInvalidParam             = 400
	ErrCodeForbid                   = 401
	ErrCodeSystemError              = 1000
	ErrCodePermissionDeny           = 1001
	ErrCodeInvalidRole              = 1002
	ErrCodeRoleUsed                 = 1003
	ErrCodeClientUsed               = 2000
	ErrCodeRoleNameDuplicate        = 2001
	ErrCodeRefreshCacheRolePolicy   = 2002
	ErrCodeAccountExist             = 2003
	ErrCodeDuplicateName            = 2004
	ErrCodeClientMediaUsedByProject = 3000
	ErrCodeUnBuried                 = 4000
	ErrCodeMediaCodeBeBusy          = 4010
	ErrCodeMediaCodeNotFound        = 4020

	ErrSchoolIdNotFound = 4101

	ErrCodeCarNameEmpty       = 5001
	ErrCodeCarMarketTimeEmpty = 5002
	ErrCodeCarCategoryIdEmpty = 5003
	ErrCodeCarPriceInvalid    = 5004
)

var ErrCodeM = map[int]string{
	ErrCodeLoginFailed:              "用户名或密码错误",
	ErrCodeInvalidParam:             "参数错误",
	ErrCodeForbid:                   "禁止访问",
	ErrCodeSystemError:              "请求出错",
	ErrCodeLoginPlease:              "请登录",
	ErrCodeInvalidRole:              "无效的角色",
	ErrCodePermissionDeny:           "无权限",
	ErrCodeRoleUsed:                 "角色已被使用",
	ErrCodeClientUsed:               "客户已被使用",
	ErrCodeRoleNameDuplicate:        "名称已被使用",
	ErrCodeRefreshCacheRolePolicy:   "权限缓存更新失败",
	ErrCodeAccountExist:             "此账号已被使用，请勿重复添加",
	ErrCodeDuplicateName:            "名称已被使用",
	ErrCodeClientMediaUsedByProject: "媒体已被项目使用",
	ErrCodeUnBuried:                 "落地页未埋点",
	ErrCodeMediaCodeBeBusy:          "改媒体码已开设投放,不可删除",
	ErrCodeMediaCodeNotFound:        "对应媒体码不存在",

	ErrSchoolIdNotFound: "对应学校不存在",

	ErrCodeCarNameEmpty:       "车型名为空",
	ErrCodeCarMarketTimeEmpty: "上市时间为空",
	ErrCodeCarCategoryIdEmpty: "车类型为空",
	ErrCodeCarPriceInvalid:    "价格区间小于等于0，或最高价格低于最低价格",
}

const (
	AdminReqId = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	RootReqId  = "ffffffff-ffff-ffff-ffff-ffffffffffff"

	UserEncryMd5 = "md5"
)

const (
	GinUserKey = "user"
	GinToken   = "Token"
	GinClient  = "Client"
	GinRoles   = "Roles"
	GinIsSuper = "IsSuper"

	HeaderUserId   = "userId"
	HeaderUsername = "username"
	HeaderSchoolId = "schoolId"
	HeaderPlatform = "platform"
	HeaderReqId    = "reqId"
	HeaderIsSuper  = "isSuper"

	GinContextKeyUserMeta = "userMeta"
)

var paramErrResponse = gin.H{
	"code":    ErrCodeInvalidParam,
	"message": ErrCodeM[ErrCodeInvalidParam],
}

var errResponse = gin.H{
	"code":    ErrCodeLoginFailed,
	"message": ErrCodeM[ErrCodeLoginFailed],
}

var sysErrResponse = gin.H{
	"code":    ErrCodeSystemError,
	"message": ErrCodeM[ErrCodeSystemError],
}

func ErrorResp(c *gin.Context, statusCode, errCode int, message string) {
	msg, ok := ErrCodeM[errCode]
	if ok && message == "" {
		message = msg
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"code": errCode, "message": message, "data": nil})
}

func SuccessResp(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": message, "data": data})
}
