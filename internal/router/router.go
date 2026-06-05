package router

import (
	"kunkun-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	// 使用自定义中间件
	// 记录请求耗时，以及状态码 并且还有是谁操作的
	router.Use(middleware.Logger())
	// 自定义全局异常恢复中间件 (防止程序 panic 导致崩溃)
	router.Use(middleware.Recovery())
	// 使用路由组
	api := router.Group("/api")
	// 使用用户路由
	SysUserRouter(api)

	// 返回路由
	return router
}
