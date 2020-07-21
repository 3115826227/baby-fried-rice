package handle

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

const (
	RootParentId = "0"
)

func SchoolDepartmentAdd(c *gin.Context) {
	var req model.ReqSchoolDepartmentAdd
	var err error
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var department model.SchoolDepartment
	department.ID = GenerateID()
	department.SchoolId = req.SchoolId
	department.Name = req.Name
	if req.ParentId == "" {
		department.ParentId = RootParentId
		department.FullName = req.Name
	} else {
		department.ParentId = req.ParentId
	}
	department.IsLeaf = 1

	tx := db.DB.Begin()
	defer func() {
		if err != nil {
			log.Logger.Warn(err.Error())
			tx.Rollback()
		}
	}()
	//如果父节点为叶子节点，修改父节点为非叶子节点
	if department.ParentId != RootParentId {
		parentDepartment, err := model.FindDepartmentById(department.ParentId)
		if parentDepartment.IsLeaf == 1 {
			var updateMap = map[string]interface{}{"is_leaf": 0}
			if err = tx.Model(model.SchoolDepartment{}).Where("id = ?", parentDepartment.ID).Update(updateMap).Error; err != nil {
				log.Logger.Warn(err.Error())
				c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
				return
			}
		}
		department.FullName = fmt.Sprintf("%v %v", parentDepartment.FullName, req.Name)
	}

	//插入当前节点
	if err = tx.Create(&department).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	err = tx.Commit().Error
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func SchoolDepartmentUpdate(c *gin.Context) {
	var req model.ReqSchoolDepartmentUpdate
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	department, err := model.FindDepartmentById(req.SchoolDepartmentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var updateMap = map[string]interface{}{"name": req.Name, "school_id": req.SchoolId, "parent_id": req.ParentId}
	if err := db.DB.Model(model.SchoolDepartment{}).Where("id = ?", department.ID).Update(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func SchoolDepartments(c *gin.Context) {
	schoolId := c.Query("school_id")
	departmentId := c.Query("department_id")

	if schoolId == "" && departmentId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	if schoolId == "" {
		department, err := model.FindDepartmentById(departmentId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
			return
		}
		schoolId = department.SchoolId
	}
	school, err := model.GetSchoolById(schoolId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	departments, err := model.GetDepartmentsBySchool(schoolId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var departmentMap = make(map[string][]model.SchoolDepartments)
	for _, department := range departments {
		if _, exist := departmentMap[department.ParentId]; !exist {
			departmentMap[department.ParentId] = make([]model.SchoolDepartments, 0)
		}
		list := departmentMap[department.ParentId]
		list = append(list, model.SchoolDepartments{
			DepartmentId:   department.ID,
			DepartmentName: department.Name,
		})
		departmentMap[department.ParentId] = list
	}

	var resp []model.SchoolDepartments
	if departmentId == "" {
		resp = Generate(departmentMap, RootParentId)
	} else {
		resp = Generate(departmentMap, departmentId)
	}

	var rsp = make([]model.RspSchoolDepartment, 0)
	rsp = append(rsp, model.RspSchoolDepartment{
		SchoolId:    school.ID,
		SchoolName:  school.Name,
		Departments: resp,
	})
	SuccessResp(c, "", rsp)
}

func SchoolDepartmentDelete(c *gin.Context) {
	departmentId := c.Query("department_id")
	department, err := model.FindDepartmentById(departmentId)
	if err != nil || department.IsLeaf != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	if err := db.DB.Delete(&department).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func Generate(departmentMap map[string][]model.SchoolDepartments, parentId string) (resp []model.SchoolDepartments) {
	resp = make([]model.SchoolDepartments, 0)
	for _, department := range departmentMap[parentId] {
		resp = append(resp, model.SchoolDepartments{
			DepartmentId:     department.DepartmentId,
			DepartmentName:   department.DepartmentName,
			ChildDepartments: Generate(departmentMap, department.DepartmentId),
		})
	}
	return
}
