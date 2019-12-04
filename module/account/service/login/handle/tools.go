package handle

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/config"
	"github.com/3115826227/baby-fried-rice/module/account/log"
	"github.com/3115826227/baby-fried-rice/module/account/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/service/model/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"strings"
	"time"
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

	HeaderUserId     = "userId"
	HeaderClientId   = "clientId"
	HeaderContractId = "contractId"
	HeaderPlatform   = "platform"
	HeaderReqId      = "reqId"
	HeaderIsSuper    = "isSuper"

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

/*
	根据用户id和创建时间生成jwt Token
*/
func GenerateToken(userID string, createTime time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     userID,
		"create_time": createTime,
	})

	return token.SignedString([]byte(config.Config.TokenSecret))
}

type UserToken struct {
	UserId     string    `json:"user_id"`
	CreateTime time.Time `json:"create_time"`
}

/*
	jwt Token解析
*/
func ExplainToken(token string) (userToken UserToken) {
	tokenInfo := strings.Split(token, ".")[1]
	var userInfoByte []byte
	userInfoByte, _ = base64.RawStdEncoding.DecodeString(tokenInfo)
	if err := json.Unmarshal(userInfoByte, &userToken); err != nil {
		fmt.Println(err.Error())
	}
	return
}

/*
	md5算法加密密码
*/
func EncodePassword(pwd string) string {
	hexStr := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	return hexStr
}

/*
	生成uuid
 */
func GenerateID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func GetToken(c *gin.Context) string {
	return c.Request.Header.Get(GinToken)
}

func GetUserMeta(c *gin.Context) *model.UserMeta {
	return c.MustGet(GinContextKeyUserMeta).(*model.UserMeta)
}

func GetRootById(id string) (root model.AccountRoot, err error) {
	err = db.DB.Find(&root).Where("id = ?", id).Error
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	return
}

func AddRoot() {
	id := GenerateID()
	now := time.Now()
	var root = model.AccountRoot{
		CommonField: model.CommonField{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		LoginName:  "root",
		Password:   EncodePassword("root"),
		Username:   "系统管理员",
		ReqId:      RootReqId,
		EncodeType: UserEncryMd5,
	}
	err := db.DB.Create(&root).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}
