package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func GetSalary(c *gin.Context)  {
	var result = make([]model.RspSalary, 0)
	var salaries = make([]model.Salary, 0)

	if err := db.DB.Find(&salaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	for _, grade := range salaries {
		result = append(result, model.RspSalary{Id: grade.ID, Name: grade.Describe})
	}

	sort.Sort(model.RspSalaries(result))

	SuccessResp(c, "", result)
}
