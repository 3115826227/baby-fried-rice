package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func GenerateUUID(c *gin.Context) {
	u, err := uuid.NewV4()
	if err != nil {
		return
	}
	c.Request.Header.Set(HeaderUUID, u.String())
	c.Set(HeaderUUID, u.String())
	c.Next()
}
