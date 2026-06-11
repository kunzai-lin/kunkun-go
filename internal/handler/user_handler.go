package handler

import (
	"kunkun-go/internal/middleware"
	"kunkun-go/internal/model"
	"kunkun-go/internal/repository"
	"kunkun-go/internal/service"
	"kunkun-go/pkg/jwt"
	"kunkun-go/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 注册用户
func RegisterUser(c *gin.Context) {
	var registerUser model.RegisterUser
	if err := c.ShouldBindJSON(&registerUser); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	user := model.RegisterUser{
		UserName:     registerUser.UserName,
		UserPassword: registerUser.UserPassword,
	}

	if err := service.CreateUser(&user); err != nil {
		response.Error(c, 500, "注册失败,用户名已存在")
		return
	}
	response.Success(c, "注册成功")

}

// 登录用户
func LoginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	var user model.RegisterUser
	// 校验用户名密码
	if err := repository.DB.Where("user_name = ? AND user_password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		response.Error(c, 401, "用户名或密码错误")
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		response.Error(c, 500, "生成 Token 失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
		},
	})
}

// 退出登录
func LogoutUser(c *gin.Context) {
	// 与 middleware.CtxUserIDKey 一致（注意是 Ctx 不是 Cxt）
	userID := c.GetUint(middleware.CtxUserIDKey)
	if userID == 0 {
		response.Error(c, 401, "用户未登录")
		return
	}
	if repository.RDB == nil {
		response.Error(c, 500, "服务未就绪")
		return
	}
	h := viper.GetInt("jwt.expire_hours")
	if h <= 0 {
		h = 72
	}
	ttl := time.Duration(h) * time.Hour
	key := "logout:" + strconv.Itoa(int(userID))
	if err := repository.RDB.Set(c.Request.Context(), key, "1", ttl).Err(); err != nil {
		response.Error(c, 500, "退出失败")
		return
	}
	response.Success(c, "退出成功")
}


func GetUserInfo(c *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: 缺少或无效的 id")
		return
	}

	user, err := service.GetUserInfo(c, req.ID)
	if err != nil {
		response.Error(c, 401, "用户不存在")
		return
	}
	response.Success(c, user)
}

// ListUsers 用户列表（分页）
func ListUsers(c *gin.Context) {
	var req struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	users, total, err := service.ListUsers(req.Page, req.PageSize)
	if err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// 更新用户
func UpdateUser(c *gin.Context) {
	var user model.SysUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := service.UpdateUser(&user); err != nil {
		response.Error(c, 400, "更新失败")
		return
	}
	response.Success(c, "更新成功")
}

// 删除用户
func DeleteUser(c *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: 缺少或无效的 id")
		return
	}

	if err := service.DeleteUser(req.ID); err != nil {
		response.Error(c, 400, "删除失败")
		return
	}
	response.Success(c, "删除成功")
}

// 创建用户
func CreateUser(c *gin.Context) {
	var user model.SysUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if err := service.CreateUser(&user); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	response.Success(c, "创建成功")
}
