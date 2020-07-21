package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUnreadMessageHandle(c *gin.Context) {
	userMeta := GetUserMeta(c)
	fmt.Println(userMeta)
}
