package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func StudentGet(c *gin.Context) {
	org := c.Query("id")
	//label := c.Query("label")
	//status := c.Query("status")
	if org == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var students = make([]model.AccountSchoolStudent, 0)
	if err := db.DB.Debug().Model(model.AccountSchoolStudent{}).Where("org_id = ?", org).Scan(&students).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var label = GetLabel(org)
	var rsp = make([]model.RspSchoolOrganizeStudent, 0)
	for _, stu := range students {
		var status = config.StudentVerifyFalse
		if stu.Status == true {
			status = config.StudentVerifyTrue
		}
		rsp = append(rsp, model.RspSchoolOrganizeStudent{
			Id:         stu.ID,
			Name:       stu.Name,
			Number:     stu.Number,
			Label:      label,
			Identify:   stu.Identify,
			Phone:      stu.Phone,
			Status:     status,
			UpdateTime: stu.UpdatedAt.Format(config.TimeLayout),
		})
	}
	SuccessResp(c, "", rsp)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}
