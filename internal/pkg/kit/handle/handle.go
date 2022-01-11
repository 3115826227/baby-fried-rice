package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	ErrCodeLoginFailed = 99
)

var ErrCodeM = map[int]string{
	ErrCodeLoginFailed:    "用户名或者密码错误",
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
}

var LoginErrResponse = gin.H{
	"code":    ErrCodeLoginFailed,
	"message": ErrCodeM[ErrCodeLoginFailed],
}

var ParamErrResponse = gin.H{
	"code":    CodeInvalidParams,
	"message": ErrCodeM[CodeInvalidParams],
}

var SysErrResponse = gin.H{
	"code":    CodeInternalError,
	"message": ErrCodeM[CodeInternalError],
}

var (
	randSource = rand.NewSource(time.Now().UnixNano())
)

func SuccessResp(c *gin.Context, message string, data interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	c.JSON(http.StatusOK, rsp.CommonResp{Code: 0, Message: message, Data: data})
}

func SuccessListResp(c *gin.Context, message string, list []interface{}, total, page, pageSize int64) {
	if list == nil {
		list = make([]interface{}, 0)
	}
	c.JSON(http.StatusOK, rsp.CommonResp{Code: 0, Message: message, Data: rsp.CommonListResp{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}})
}

func SystemErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, SysErrResponse)
}

func FailedResp(c *gin.Context, code int) {
	msg, ok := ErrCodeM[code]
	if !ok {
		c.JSON(http.StatusOK, SysErrResponse)
		return
	}
	c.JSON(http.StatusOK, rsp.CommonResp{Code: code, Message: msg, Data: make(map[string]interface{})})
}

func ErrorResp(c *gin.Context, statusCode, errCode int, message string) {
	msg, ok := ErrCodeM[errCode]
	if ok && message == "" {
		message = msg
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"code": errCode, "message": message, "data": nil})
}

func EncodePassword(id, pwd string) string {
	hexStr := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v:%v", id, pwd))))
	return hexStr
}

func GenerateID() string {
	return uuid.NewV4().String()
}

func GeneratePhoneCode() string {
	return fmt.Sprintf("%04v", rand.New(randSource).Int31n(10000))
}

//生成八位数字
func GenerateSerialNumber() string {
	return fmt.Sprintf("1%08v", rand.New(randSource).Int31n(1000000))
}

//生成十二位数字
func GenerateSerialNumberByLen(len int) string {
	var res = 1
	for i := 0; i < len; i++ {
		res *= 10
	}
	return fmt.Sprintf("1%0"+strconv.Itoa(len)+"v", rand.New(randSource).Int31n(int32(res)))
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

func PageHandle(c *gin.Context) (req requests.PageCommonReq, err error) {
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = fmt.Sprintf("%v", constant.DefaultPage)
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return
	}
	if page <= 0 {
		page = constant.DefaultPage
	}
	pageSizeStr := c.Query("page_size")
	if pageSizeStr == "" {
		pageSizeStr = fmt.Sprintf("%v", constant.DefaultPageSize)
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return
	}
	if pageSize <= 0 {
		pageSize = constant.DefaultPageSize
	}
	req.Page = int64(page)
	req.PageSize = int64(pageSize)
	return
}
