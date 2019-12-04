package handle

import (
	"github.com/gin-gonic/gin"
)

const (
	HeaderUserId = "userId"

	StrangerID = ""
)

func HasUser(c *gin.Context) (string, bool) {
	userID := c.GetHeader(HeaderUserId)
	if userID == "" {
		return userID, false
	}
	return userID, true
}
