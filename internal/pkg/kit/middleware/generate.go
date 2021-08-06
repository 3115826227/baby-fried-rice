package middleware

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func GenerateUUID(c *gin.Context) {
	u := uuid.NewV4().String()
	c.Request.Header.Set(handle.HeaderUUID, u)
	c.Set(handle.HeaderUUID, u)
	c.Next()
}
