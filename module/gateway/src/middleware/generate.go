package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func GenerateUUID(c *gin.Context) {
	u := uuid.NewV4()
	c.Request.Header.Set(HeaderUUID, u.String())
	c.Set(HeaderUUID, u.String())
	c.Next()
}
