package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminLoginLogAdd(id, ip string, now time.Time) {
	var count int
	if err := db.DB.Debug().Model(&model.AccountAdminLoginLog{}).Where("admin_id = ?", id).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var loginLog = &model.AccountAdminLoginLog{
		AdminID:    id,
		IP:         ip,
		LoginCount: count + 1,
		LoginTime:  now,
	}
	if err := db.DB.Debug().Create(&loginLog).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	go UpdateIp(ip)
}

func AdminLoginLog(c *gin.Context) {
	id := c.Query("id")
	var logs = make([]model.AccountAdminLoginLog, 0)
	if id != "" {
		logs = model.GetAdminLoginLog(id)
	} else {
		logs = model.GetAdminLoginLog()
	}
	var resp = make([]model.RspAdminLoginLog, 0)
	var ipMap = make(map[string]model.Ip)
	var ips = make([]string, 0)
	var adminMap = make(map[string]model.AccountAdmin)
	var ids = make([]string, 0)
	for _, l := range logs {
		if _, exist := adminMap[l.AdminID]; !exist {
			ids = append(ids, l.AdminID)
			adminMap[l.AdminID] = model.AccountAdmin{}
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
	var schoolIds = make([]string, 0)
	var schoolMap = make(map[string]model.School)
	admins := model.GetAdmin(ids...)
	for _, a := range admins {
		adminMap[a.ID] = a
		if _, exist := schoolMap[a.SchoolId]; !exist {
			schoolIds = append(schoolIds, a.SchoolId)
			schoolMap[a.SchoolId] = model.School{}
		}
	}
	schools := model.GetSchool(schoolIds...)
	for _, s := range schools {
		schoolMap[s.ID] = s
	}
	for _, l := range logs {
		var loginLog model.RspAdminLoginLog
		loginLog.AdminId = l.AdminID
		loginLog.LoginName = adminMap[l.AdminID].LoginName
		loginLog.Username = adminMap[l.AdminID].Username
		loginLog.Phone = adminMap[l.AdminID].Phone
		loginLog.Count = l.LoginCount
		loginLog.Ip = l.IP
		loginLog.Area = ipMap[l.IP].Describe
		loginLog.School = schoolMap[adminMap[l.AdminID].SchoolId].Name
		loginLog.Time = l.LoginTime.Format(config.TimeLayout)
		resp = append(resp, loginLog)
	}

	SuccessResp(c, "", resp)
}
