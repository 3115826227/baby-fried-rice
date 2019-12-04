package register

import "github.com/gin-gonic/gin"

/*
	服务注册
*/
func Route(apiGroup *gin.RouterGroup) {

	apiGroup.GET("/register", Register)
}
