package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

func StudentGet(c *gin.Context) {
	organize := c.Query("organize")
	if organize == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var students = make([]model.AccountSchoolStudent, 0)
	if err := db.DB.Debug().Where("org_id = ?", organize).Find(&students).Error; err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var rspStudents = make([]model.RspSchoolStudent, 0)
	for _, s := range students {
		rspStudents = append(rspStudents, model.RspSchoolStudent{
			Id:         s.ID,
			Name:       s.Name,
			Identify:   s.Identify,
			Status:     s.Status,
			Number:     s.Number,
			Phone:      s.Phone,
			OrgId:      s.OrgId,
			CreateTime: s.CreatedAt.Format(config.TimeLayout),
			UpdateTime: s.UpdatedAt.Format(config.TimeLayout),
		})
	}
	sort.Sort(model.RspSchoolStudents(rspStudents))
	SuccessResp(c, "", rspStudents)
}

func StudentAdd(c *gin.Context) {
	var req model.ReqSchoolStudentAdd
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var student = model.AccountSchoolStudent{
		Name:     req.Name,
		Identify: req.Identify,
		Status:   false,
		Number:   req.Number,
		Phone:    "",
		OrgId:    req.Organize,
	}
	var now = time.Now()
	student.ID = GenerateID()
	student.CreatedAt = now
	student.UpdatedAt = now
	if err := db.DB.Debug().Create(&student).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}
