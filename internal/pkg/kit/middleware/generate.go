package middleware

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func GenerateUUID(c *gin.Context) {
	u := uuid.NewV4()
	c.Request.Header.Set(handle.HeaderUUID, u.String())
	c.Set(handle.HeaderUUID, u.String())
	c.Next()
}
