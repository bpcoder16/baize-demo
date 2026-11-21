package router

import (
	"github.com/bpcoder16/baize-demo/internal/http/handler"
	"github.com/gin-gonic/gin"
)

func registerAppRoutes(app *gin.RouterGroup) {
	app.GET("/ping", handler.Ping)
}
