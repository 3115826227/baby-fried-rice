package handle

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	HeaderUserId   = "userId"
	HeaderUsername = "username"
	HeaderSchoolId = "schoolId"
	HeaderPlatform = "platform"
	HeaderReqId    = "reqId"
	HeaderIsSuper  = "isSuper"

	HeaderIP     = "IP"
	HeaderArea   = "Area"
	HeaderBrower = "Brower"
	HeaderOS     = "OS"

	GinContextKeyUserMeta = "userMeta"
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
	ErrCodeRefreshCacheRolePolicy: "权限缓存更新失败",
	ErrCodeAccountExist:           "此账号已被使用，请勿重复添加",
	ErrCodeAccountNotFound:        "对应账号不存在",
	ErrCodeDuplicateName:          "名称已被使用",
}

const (
	AdminPassword = "1234"
	UserEncryMd5  = "md5"
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

func LoginPost(url string, payload []byte, header http.Header) (data []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	req.Header = header
	req.Header.Add("Content-Type", "application/json")
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	return
}

func Post(url string, payload []byte) (data []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
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
	id := uuid.NewV4()
	return id.String()
}

func GetUserMeta(c *gin.Context) *model.UserMeta {
	return c.MustGet(GinContextKeyUserMeta).(*model.UserMeta)
}
