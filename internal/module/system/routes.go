package system

import (
	"echo-framework/internal/module/system/user"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes 注册系统模块路由
func RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/sys")
	user.RegisterRoutes(api.Group("/users"))
}
