package service

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"github.com/3115826227/baby-fried-rice/module/crawler/model/db"
	"github.com/3115826227/baby-fried-rice/module/crawler/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
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

const (
	TrainMetaField  = "train:meta:field"
	TrainMetaMember = "train:meta:member"
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

func GetStation() (stations []model.Station) {
	if err := db.DB.Find(&stations).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	sort.Sort(model.Stations(stations))
	return
}

func IsCrawlTrainInDate(date, train string) bool {
	return redis.HashExist(fmt.Sprintf("%v:%v", TrainMetaField, date), train)
}

func IsValidDate(date string) bool {
	_, err := time.Parse(config.DayLayout, date)
	if err != nil {
		return false
	}
	return true
}

type TrainComputeInfo struct {
	Train    model.TrainMeta
	Stations []model.TrainStationRelation
}

func (info *TrainComputeInfo) ComputeRunningAndOverDay() {
	for i := 1; i < len(info.Stations); i++ {
		info.Stations[i].OverDay = info.Stations[i-1].OverDay
		/*
			比较本站到达时间与上一站开出时间，判断是否多出了一天
		*/
		if isOverDay(info.Stations[i].ArriveTime, info.Stations[i-1].StartTime) {
			info.Stations[i].OverDay += 1
		}
		/*
			比较本站开出时间与本站到达时间，判断是否多出了一天
		*/
		if isOverDay(info.Stations[i].StartTime, info.Stations[i].ArriveTime) {
			info.Stations[i].OverDay += 1
		}
	}
	info.Train.OverDay = info.Stations[len(info.Stations)-1].OverDay
	info.Train.RunningMinute = computeRunning(info.Train.StartTime, info.Train.ArriveTime, info.Train.OverDay)
	return
}

func computeRunning(startTime, arriveTime string, overDay int) int {
	return HourMinuteToConvert(arriveTime) - HourMinuteToConvert(startTime) + overDay*24*60
}

func isOverDay(now, last string) bool {
	return HourMinuteToConvert(now) < HourMinuteToConvert(last)
}

func HourMinuteToConvert(now string) int {
	nowSlice := strings.Split(now, ":")
	hour, _ := strconv.Atoi(nowSlice[0])
	minute, _ := strconv.Atoi(nowSlice[1])
	return hour*60 + minute
}
