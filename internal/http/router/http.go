// Package router 注册 Demo 的 HTTP 路由。
package router

import "github.com/gin-gonic/gin"

// Register 将 App 和 Admin 路由挂载到各自的路径前缀。
func Register(engine *gin.Engine) {
	registerAppRoutes(engine.Group("/app"))
	registerAdminRoutes(engine.Group("/admin"))
}
