package router

import (
	"github.com/bpcoder16/baize-demo/internal/http/handler"
	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(admin *gin.RouterGroup) {
	admin.GET("/ping", handler.Ping)
}
