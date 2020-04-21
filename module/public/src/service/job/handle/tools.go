package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/log"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
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

func AddSubject() {
	var beans = make([]interface{}, 0)
	beans = append(beans, &model.Subject{Name: "语文"})
	beans = append(beans, &model.Subject{Name: "数学"})
	beans = append(beans, &model.Subject{Name: "外语"})
	beans = append(beans, &model.Subject{Name: "物理"})
	beans = append(beans, &model.Subject{Name: "生物"})
	beans = append(beans, &model.Subject{Name: "化学"})
	beans = append(beans, &model.Subject{Name: "政治"})
	beans = append(beans, &model.Subject{Name: "历史"})
	beans = append(beans, &model.Subject{Name: "地理"})

	err := db.CreateMulti(beans...)
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

func AddGrade() {
	var beans = make([]interface{}, 0)
	beans = append(beans, &model.Grade{Name: "一年级"})
	beans = append(beans, &model.Grade{Name: "二年级"})
	beans = append(beans, &model.Grade{Name: "三年级"})
	beans = append(beans, &model.Grade{Name: "四年级"})
	beans = append(beans, &model.Grade{Name: "五年级"})
	beans = append(beans, &model.Grade{Name: "六年级"})
	beans = append(beans, &model.Grade{Name: "初一"})
	beans = append(beans, &model.Grade{Name: "初二"})
	beans = append(beans, &model.Grade{Name: "初三"})
	beans = append(beans, &model.Grade{Name: "高一"})
	beans = append(beans, &model.Grade{Name: "高二"})
	beans = append(beans, &model.Grade{Name: "高三"})

	err := db.CreateMulti(beans...)
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

func addCourse(gradeId, subjectId int) *model.Course {
	var course = new(model.Course)
	var grade = &model.Grade{ID: gradeId}
	var subject = &model.Subject{ID: subjectId}
	db.DB.First(&grade)
	db.DB.First(&subject).Where("id = ?", subjectId)
	course.SubjectId = subjectId
	course.GradeId = gradeId
	course.Name = grade.Name + subject.Name
	return course
}

func AddCourse() {
	var beans = make([]interface{}, 0)

	beans = append(beans, addCourse(1, 1))
	beans = append(beans, addCourse(1, 2))
	beans = append(beans, addCourse(1, 3))

	beans = append(beans, addCourse(2, 1))
	beans = append(beans, addCourse(2, 2))
	beans = append(beans, addCourse(2, 3))

	beans = append(beans, addCourse(3, 1))
	beans = append(beans, addCourse(3, 2))
	beans = append(beans, addCourse(3, 3))

	beans = append(beans, addCourse(4, 1))
	beans = append(beans, addCourse(4, 2))
	beans = append(beans, addCourse(4, 3))

	beans = append(beans, addCourse(5, 1))
	beans = append(beans, addCourse(5, 2))
	beans = append(beans, addCourse(5, 3))

	beans = append(beans, addCourse(6, 1))
	beans = append(beans, addCourse(6, 2))
	beans = append(beans, addCourse(6, 3))

	beans = append(beans, addCourse(7, 1))
	beans = append(beans, addCourse(7, 2))
	beans = append(beans, addCourse(7, 3))
	beans = append(beans, addCourse(7, 5))
	beans = append(beans, addCourse(7, 7))
	beans = append(beans, addCourse(7, 8))
	beans = append(beans, addCourse(7, 9))

	beans = append(beans, addCourse(8, 1))
	beans = append(beans, addCourse(8, 2))
	beans = append(beans, addCourse(8, 3))
	beans = append(beans, addCourse(8, 4))
	beans = append(beans, addCourse(8, 5))
	beans = append(beans, addCourse(8, 7))
	beans = append(beans, addCourse(8, 8))
	beans = append(beans, addCourse(8, 9))

	beans = append(beans, addCourse(9, 1))
	beans = append(beans, addCourse(9, 2))
	beans = append(beans, addCourse(9, 3))
	beans = append(beans, addCourse(9, 4))
	beans = append(beans, addCourse(9, 6))
	beans = append(beans, addCourse(9, 7))
	beans = append(beans, addCourse(9, 8))

	beans = append(beans, addCourse(10, 1))
	beans = append(beans, addCourse(10, 2))
	beans = append(beans, addCourse(10, 3))
	beans = append(beans, addCourse(10, 4))
	beans = append(beans, addCourse(10, 5))
	beans = append(beans, addCourse(10, 6))
	beans = append(beans, addCourse(10, 7))
	beans = append(beans, addCourse(10, 8))
	beans = append(beans, addCourse(10, 9))

	beans = append(beans, addCourse(11, 1))
	beans = append(beans, addCourse(11, 2))
	beans = append(beans, addCourse(11, 3))
	beans = append(beans, addCourse(11, 4))
	beans = append(beans, addCourse(11, 5))
	beans = append(beans, addCourse(11, 6))
	beans = append(beans, addCourse(11, 7))
	beans = append(beans, addCourse(11, 8))
	beans = append(beans, addCourse(11, 9))

	beans = append(beans, addCourse(12, 1))
	beans = append(beans, addCourse(12, 2))
	beans = append(beans, addCourse(12, 3))
	beans = append(beans, addCourse(12, 4))
	beans = append(beans, addCourse(12, 5))
	beans = append(beans, addCourse(12, 6))
	beans = append(beans, addCourse(12, 7))
	beans = append(beans, addCourse(12, 8))
	beans = append(beans, addCourse(12, 9))

	err := db.CreateMulti(beans...)
	if err != nil {
		log.Logger.Warn(err.Error())
	}

}

func AddSalary() {
	var beans = make([]interface{}, 0)
	beans = append(beans, &model.Salary{Describe: "50以下", Min: 0, Max: 50})
	beans = append(beans, &model.Salary{Describe: "50 ~ 60", Min: 50, Max: 60})
	beans = append(beans, &model.Salary{Describe: "60 ~ 70", Min: 60, Max: 70})
	beans = append(beans, &model.Salary{Describe: "70 ~ 80", Min: 70, Max: 80})
	beans = append(beans, &model.Salary{Describe: "80以上", Min: 80, Max: 200})

	err := db.CreateMulti(beans...)
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}
