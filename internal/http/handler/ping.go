// Package handler 实现 Demo 的 HTTP 处理函数，并维护各接口的 Swag 注释。
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary App 连通性检查
// @Description 仅检查 HTTP 服务是否可达，不检查数据库等依赖。
// @x-module-app true
// @Tags 连通性检查
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /app/ping [get]

// @Summary Admin 连通性检查
// @Description 仅检查 HTTP 服务是否可达，不检查数据库等依赖。
// @x-module-admin true
// @Tags 连通性检查
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /admin/ping [get]

// Ping 返回纯文本 pong，供 App 和 Admin 检查 HTTP 连通性。
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
