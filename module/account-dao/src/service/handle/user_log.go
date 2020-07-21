package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"time"
)

func UserLoginLogAdd(id, ip string, now time.Time) {
	var count int
	if err := db.DB.Debug().Model(&model.AccountUserLoginLog{}).Where("user_id = ?", id).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var loginLog = &model.AccountUserLoginLog{
		UserID:     id,
		IP:         ip,
		LoginCount: count + 1,
		LoginTime:  now,
	}
	if err := db.DB.Debug().Create(&loginLog).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	go UpdateIp(ip)
}

func UserLoginLog(c *gin.Context) {
	id := c.Query("id")
	var logs = make([]model.AccountUserLoginLog, 0)
	if id != "" {
		logs = model.GetUserLoginLog(id)
	} else {
		logs = model.GetUserLoginLog()
	}
	var resp = make([]model.RspUserLoginLog, 0)
	var ipMap = make(map[string]model.Ip)
	var ips = make([]string, 0)
	var detailMap = make(map[string]model.AccountUserDetail)
	var userMap = make(map[string]model.AccountUser)
	var ids = make([]string, 0)
	for _, l := range logs {
		if _, exist := detailMap[l.UserID]; !exist {
			ids = append(ids, l.UserID)
			detailMap[l.UserID] = model.AccountUserDetail{}
		}
		if _, exist := ipMap[l.IP]; !exist {
			ips = append(ips, l.IP)
			ipMap[l.IP] = model.Ip{}
		}
	}
	var ipList = make([]model.Ip, 0)
	if err := db.DB.Debug().Where("ip in (?)", ips).Find(&ipList).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	for _, item := range ipList {
		ipMap[item.Ip] = item
	}
	users := model.GetUsers(ids)
	details := model.GetUserDetail(ids...)
	for _, d := range details {
		detailMap[d.ID] = d
	}
	for _, u := range users {
		userMap[u.ID] = u
	}
	for _, l := range logs {
		var loginLog model.RspUserLoginLog
		loginLog.UserId = l.UserID
		loginLog.LoginName = userMap[l.UserID].LoginName
		loginLog.Username = detailMap[l.UserID].Username
		loginLog.Phone = detailMap[l.UserID].Phone
		loginLog.Count = l.LoginCount
		loginLog.Ip = l.IP
		loginLog.Area = ipMap[l.IP].Describe
		loginLog.Time = l.LoginTime.Format(config.TimeLayout)
		resp = append(resp, loginLog)
	}

	SuccessResp(c, "", resp)
}
