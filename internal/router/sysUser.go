package router

import (
	"kunkun-go/internal/handler"
	"kunkun-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SysUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	// 注册、登录不需要 token
	user.POST("/register", handler.RegisterUser)
	user.POST("/login", handler.LoginUser)

	// 其余接口需要携带有效 JWT（Header: Authorization: Bearer <token>）
	protected := user.Group("")
	protected.Use(middleware.JWTAuth())
	protected.GET("/info", handler.GetUserInfo)
	protected.PUT("/update", handler.UpdateUser)
	protected.DELETE("/delete", handler.DeleteUser)
	protected.POST("/create", handler.CreateUser)
}
