package handle

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"time"
)

const (
	ErrCodeLoginFailed  = 99
	ErrCodeInvalidParam = 400
	ErrCodeSystemError  = 1000
)

var ErrCodeM = map[int]string{
	ErrCodeLoginFailed:  "用户名或密码错误",
	ErrCodeInvalidParam: "参数错误",
	ErrCodeSystemError:  "请求出错",
}

var LoginErrResponse = gin.H{
	"code":    ErrCodeLoginFailed,
	"message": ErrCodeM[ErrCodeLoginFailed],
}

var ParamErrResponse = gin.H{
	"code":    ErrCodeInvalidParam,
	"message": ErrCodeM[ErrCodeInvalidParam],
}

var SysErrResponse = gin.H{
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

func EncodePassword(pwd string) string {
	hexStr := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	return hexStr
}

func GenerateID() string {
	return uuid.NewV4().String()
}

//生成八位数字
func GenerateSerialNumber() string {
	return fmt.Sprintf("1%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

/*
	根据用户id和创建时间生成jwt Token
*/
func GenerateToken(userID string, createTime time.Time, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     userID,
		"create_time": createTime,
	})

	return token.SignedString([]byte(tokenSecret))
}

func ResponseHandle(data []byte) (ok bool, err error) {
	var resp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	ok = resp.Code == SuccessCode
	return
}
