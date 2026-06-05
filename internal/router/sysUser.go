package router

import (
	"kunkun-go/internal/handler"

	"github.com/gin-gonic/gin"
)

func SysUserRouter(router *gin.RouterGroup) {
	auth := router.Group("/user")
	auth.POST("/register", handler.RegisterUser)
	auth.POST("/login", handler.LoginUser)
	auth.GET("/info", handler.GetUserInfo)
	auth.PUT("/update", handler.UpdateUser)
	auth.DELETE("/delete", handler.DeleteUser)
	auth.POST("/create", handler.CreateUser)
}
