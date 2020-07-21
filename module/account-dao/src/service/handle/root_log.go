package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"time"
)

func RootLoginLogAdd(id, ip string, now time.Time) {
	var count int
	if err := db.DB.Debug().Model(&model.AccountRootLoginLog{}).Where("root_id = ?", id).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var loginLog = &model.AccountRootLoginLog{
		RootID:     id,
		IP:         ip,
		LoginCount: count + 1,
		LoginTime:  now,
	}
	if err := db.DB.Debug().Create(&loginLog).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	go UpdateIp(ip)
}

func RootLoginLog(c *gin.Context) {
	id := c.Query("id")
	var logs = make([]model.AccountRootLoginLog, 0)
	if id != "" {
		logs = model.GetRootLoginLog(id)
	} else {
		logs = model.GetRootLoginLog()
	}
	var resp = make([]model.RspRootLoginLog, 0)
	var ipMap = make(map[string]model.Ip)
	var ips = make([]string, 0)
	var rootMap = make(map[string]model.AccountRoot)
	var ids = make([]string, 0)
	for _, l := range logs {
		if _, exist := rootMap[l.RootID]; !exist {
			ids = append(ids, l.RootID)
			rootMap[l.RootID] = model.AccountRoot{}
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
	roots := model.GetRoots(ids)
	for _, r := range roots {
		rootMap[r.ID] = r
	}
	for _, l := range logs {
		var loginLog model.RspRootLoginLog
		loginLog.RootId = l.RootID
		loginLog.LoginName = rootMap[l.RootID].LoginName
		loginLog.Username = rootMap[l.RootID].Username
		loginLog.Phone = rootMap[l.RootID].Phone
		loginLog.Count = l.LoginCount
		loginLog.Ip = l.IP
		loginLog.Area = ipMap[l.IP].Describe
		loginLog.Time = l.LoginTime.Format(config.TimeLayout)
		resp = append(resp, loginLog)
	}

	SuccessResp(c, "", resp)
}
