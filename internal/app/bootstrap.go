package app

import (
	config2 "echo-framework/internal/config"
	"echo-framework/internal/infrastructure/db"
	"echo-framework/internal/infrastructure/log"
	"echo-framework/internal/infrastructure/monitor"
	customMiddleware "echo-framework/internal/middleware"
	"echo-framework/internal/module/system"

	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start 启动应用程序
func Start() {
	conf, err := config2.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load configuration file: %v", err)
		return
	}

	err = setupInfrastructure(conf)
	if err != nil {
		fmt.Printf("failed to initialize services: %v", err)
		return
	}

	e := echo.New()
	registerMiddleware(e)
	registerRoutes(e, conf)

	err = e.Start(":" + conf.App.Port)
	if err != nil {
		return
	}
}

// setupInfrastructure 初始化基础组件
func setupInfrastructure(conf *config2.Config) (err error) {
	// 初始化日志
	err = log.InitLogger(conf)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}
	// 初始化数据库
	if conf.Database.Sync {
		err = db.SyncDatabase()
		if err != nil {
			return fmt.Errorf("failed to sync database: %v", err)
		}
	}
	return nil
}

// registerMiddleware 注册中间件
func registerMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(customMiddleware.Logger())
	e.Use(customMiddleware.Jwt())
}

// registerRoutes 注册路由
func registerRoutes(e *echo.Echo, conf *config2.Config) {
	// 注册 statsviz 运行时监控
	if conf.Monitor.Enable {
		monitor.RegisterRoutes(e, conf.Monitor.Path)
	}
	// 注册系统模块路由
	system.RegisterRoutes(e)
}
