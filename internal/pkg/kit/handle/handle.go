package handle

import (
	"crypto/md5"
	"fmt"
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

var loginErrResponse = gin.H{
	"code":    ErrCodeLoginFailed,
	"message": ErrCodeM[ErrCodeLoginFailed],
}

var paramErrResponse = gin.H{
	"code":    ErrCodeInvalidParam,
	"message": ErrCodeM[ErrCodeInvalidParam],
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
