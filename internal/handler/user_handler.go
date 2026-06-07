package handler

import (
	"kunkun-go/internal/model"
	"kunkun-go/internal/repository"
	"kunkun-go/internal/service"
	"kunkun-go/pkg/jwt"
	"kunkun-go/pkg/response"

	"github.com/gin-gonic/gin"
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

func GetUserInfo(c *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: 缺少或无效的 id")
		return
	}

	user, err := service.GetUserInfo(req.ID)
	if err != nil {
		response.Error(c, 401, "用户不存在")
		return
	}
	response.Success(c, user)
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
