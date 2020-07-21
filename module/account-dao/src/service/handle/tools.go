package handle

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
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

func IsDuplicateLoginNameByUser(loginName string) bool {
	var users = make([]model.AccountUser, 0)
	var count = 0
	if err := db.DB.Where("login_name = ?", loginName).Find(&users).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return true
	}
	return count != 0
}

func GenerateID() string {
	return uuid.NewV4().String()
}

//生成八位数字
func GenerateSerialNumber() string {
	return fmt.Sprintf("1%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
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

func UpdateIp(ip string) (describe string) {
	url := fmt.Sprintf("%v/%v?lang=zh-CN", config.GetIpAddress, ip)
	data, err := Get(url)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var resp model.RespIp
	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	describe = resp.RegionName + resp.City
	var accountIp = &model.Ip{
		Ip:         ip,
		Province:   resp.RegionName,
		City:       resp.City,
		Describe:   resp.RegionName + resp.City,
		UpdateTime: time.Now(),
	}
	if err := db.DB.Debug().Model(&model.Ip{}).Save(&accountIp).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
