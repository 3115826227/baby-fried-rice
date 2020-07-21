package handle

import (
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/redis"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

const (
	CodeInvalidParams     = 400
	CodeRequiredLogin     = 401
	CodeRequiredForbidden = 403
	CodeNotFound          = 404
	CodeInternalError     = 500
	CodeServiceNotFound   = 502
)

const (
	CodeInvalidParamsMsg     = "参数错误"
	CodeRequiredLoginMsg     = "请登录"
	CodeRequiredForbiddenMsg = "权限不够"
	CodeNotFoundMsg          = "未找到服务"
	CodeInternalErrorMsg     = "服务器错误"
	CodeServiceNotFoundMsg   = "服务不存在"
)

const (
	ErrCodeLoginFailed            = 99
	ErrCodeLoginPlease            = 100
	ErrCodeUserBanned             = 101
	ErrCodeInvalidParam           = 400
	ErrCodeForbid                 = 401
	ErrCodeSystemError            = 1000
	ErrCodePermissionDeny         = 1001
	ErrCodeInvalidRole            = 1002
	ErrCodeRoleUsed               = 1003
	ErrCodeClientUsed             = 2000
	ErrCodeRoleNameDuplicate      = 2001
	ErrCodeRefreshCacheRolePolicy = 2002
	ErrCodeAccountExist           = 2003
	ErrCodeDuplicateName          = 2004

	ErrSchoolIdNotFound    = 4101
	ErrCodeAccountNotFound = 4102

	ErrCodeCarNameEmpty       = 5001
	ErrCodeCarMarketTimeEmpty = 5002
	ErrCodeCarCategoryIdEmpty = 5003
	ErrCodeCarPriceInvalid    = 5004
)

var ErrCodeM = map[int]string{
	ErrCodeLoginFailed:            "用户名或密码错误",
	ErrCodeInvalidParam:           "参数错误",
	ErrCodeForbid:                 "禁止访问",
	ErrCodeSystemError:            "请求出错",
	ErrCodeLoginPlease:            "请登录",
	ErrCodeInvalidRole:            "无效的角色",
	ErrCodePermissionDeny:         "无权限",
	ErrCodeRoleUsed:               "角色已被使用",
	ErrCodeClientUsed:             "客户已被使用",
	ErrCodeRoleNameDuplicate:      "名称已被使用",
	ErrCodeRefreshCacheRolePolicy: "权限缓存更新失败",
	ErrCodeAccountExist:           "此账号已被使用，请勿重复添加",
	ErrCodeAccountNotFound:        "对应账号不存在",
	ErrCodeDuplicateName:          "名称已被使用",

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

const (
	TokenPrefix = "token"
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

func SuccessResp(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": message, "data": data})
}

func ErrorResp(c *gin.Context, statusCode, errCode int, message string) {
	msg, ok := ErrCodeM[errCode]
	if ok && message == "" {
		message = msg
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"code": errCode, "message": message, "data": nil})
}

func GenerateID() string {
	return uuid.NewV4().String()
}

func GetUserMeta(c *gin.Context) *model.UserMeta {
	return c.MustGet(GinContextKeyUserMeta).(*model.UserMeta)
}

func GetUserMetaByToken(c *gin.Context, token string) (userMeta *model.UserMeta, err error) {
	if token == "" {
		ErrorResp(c, http.StatusUnauthorized, CodeRequiredLogin, CodeRequiredLoginMsg)
		c.Abort()
		return
	}

	tokenKey := fmt.Sprintf("%v:%v", TokenPrefix, token)
	var str string
	str, err = redis.Get(tokenKey)
	if err != nil {
		log.Logger.Warn(err.Error())
		ErrorResp(c, http.StatusUnauthorized, CodeRequiredLogin, CodeRequiredLoginMsg)
		c.Abort()
		return
	}
	err = json.Unmarshal([]byte(str), &userMeta)
	if err != nil {
		ErrorResp(c, http.StatusUnauthorized, CodeRequiredLogin, CodeRequiredLoginMsg)
		c.Abort()
		return
	}
	return
}

func Get(url string) (data []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var res *http.Response
	res, err = client.Do(req)
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	return
}
