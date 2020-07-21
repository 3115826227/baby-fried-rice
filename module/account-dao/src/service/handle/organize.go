package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

func OrganizeGetHandle(c *gin.Context) {
	schoolId := c.Query("school_id")
	label := c.Query("label")
	status := c.Query("status")

	var rsp = make([]model.RspSchoolOrganize, 0)
	var organizes = make([]model.AccountSchoolOrganize, 0)
	var sql = fmt.Sprintf(`select a.*, b.count from account_school_organize as a 
left join (
select org_id, count(*) as count  from account_school_student
GROUP BY org_id ) as b
on a.id = b.org_id
where a.school_id = '%v'
`, schoolId)
	if label != "" {
		sql = fmt.Sprintf(`%v and a.label = '%v'`, sql, label)
	}
	if status != "" && status != "all" {
		sql = fmt.Sprintf(`%v and a.status = '%v'`, sql, status)
	}
	if err := db.DB.Debug().Raw(sql).Scan(&organizes).Error; err != nil {
		log.Logger.Warn(err.Error())
		SuccessResp(c, "", rsp)
		return
	}

	var mp = make(map[string]model.RspSchoolOrganize)
	for _, org := range organizes {
		mp[org.ID] = model.RspSchoolOrganize{
			Id:         org.ID,
			Label:      org.Label,
			ParentId:   org.ParentId,
			Status:     org.Status,
			SchoolId:   org.SchoolId,
			CreateTime: org.CreatedAt.Format(config.TimeLayout),
			UpdateTime: org.UpdatedAt.Format(config.TimeLayout),
			Children:   make([]model.RspSchoolOrganize, 0),
		}
	}

	for _, org := range mp {
		if org.ParentId == config.RootSchoolOrganizeId {
			DFSGetSchoolOrganize(&org, mp)
			rsp = append(rsp, org)
		}
	}

	sort.Sort(model.RspSchoolOrganizes(rsp))

	SuccessResp(c, "", rsp)
}

func DFSGetSchoolOrganize(organize *model.RspSchoolOrganize, mp map[string]model.RspSchoolOrganize) {
	for _, org := range mp {
		if org.ParentId == organize.Id {
			DFSGetSchoolOrganize(&org, mp)
			organize.Count += org.Count
			organize.Children = append(organize.Children, org)
		}
	}
	sort.Sort(model.RspSchoolOrganizes(organize.Children))
}

func OrganizeExistHandle(c *gin.Context) {
	id := c.Query("id")

	var count int
	if err := db.DB.Debug().Model(&model.AccountSchoolOrganize{}).Where("id = ?", id).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if count != 1 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	SuccessResp(c, "", nil)
}

func OrganizeAdd(c *gin.Context) {
	var req model.ReqSchoolOrganizeAdd
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var org = model.AccountSchoolOrganize{
		Label:    req.Label,
		ParentId: req.ParentId,
		SchoolId: req.SchoolId,
		Status:   req.Status,
	}
	var now = time.Now()
	org.ID = GenerateID()
	org.CreatedAt = now
	org.UpdatedAt = now
	if err := db.DB.Debug().Create(&org).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func OrganizeUpdate(c *gin.Context) {
	var req model.ReqSchoolOrganizedUpdate
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if err := db.DB.Debug().Model(&model.AccountSchoolOrganize{}).Where("id = ?", req.Id).Update("label", req.Label).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func OrganizeStatus(c *gin.Context) {
	var req model.ReqSchoolOrganizedStatusUpdate
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if err := db.DB.Debug().Model(&model.AccountSchoolOrganize{}).Where("id = ?", req.Id).Update("status", req.Status).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func GetLabel(org string) (label string) {
	var organizes = make([]model.AccountSchoolOrganize, 0)
	if err := db.DB.Model(&model.AccountSchoolOrganize{}).Find(&organizes).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	var mp = make(map[string]model.AccountSchoolOrganize)
	for _, o := range organizes {
		mp[o.ID] = o
	}
	return DFSGetLabel(mp[org].Label, mp[mp[org].ParentId], mp)
}

func DFSGetLabel(label string, org model.AccountSchoolOrganize, mp map[string]model.AccountSchoolOrganize) string {
	label = fmt.Sprintf("%v-%v", org.Label, label)
	if org.ParentId == config.RootSchoolOrganizeId {
		return label
	}
	return DFSGetLabel(label, mp[org.ParentId], mp)
}

func OrganizeDelete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if err := db.DB.Debug().Where("id = ?", id).Delete(&model.AccountSchoolOrganize{}).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}
