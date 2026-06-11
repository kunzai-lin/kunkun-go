package service

import (
	"encoding/json"
	"fmt"
	"kunkun-go/internal/model"
	"kunkun-go/internal/repository"

	"github.com/gin-gonic/gin"
)

// 获取用户信息
func GetUserInfo(ctx *gin.Context, id uint) (model.SysUser, error) {

	// 从 Redis 中获取用户信息
	key := fmt.Sprintf("user:info:%d", id)
	cached, err := repository.RDB.Get(ctx, key).Result()
	if err == nil {
		var user model.SysUser
		json.Unmarshal([]byte(cached), &user)
		return user, nil
	}
	// 如果 Redis 中没有用户信息，则从数据库中获取
	var user model.SysUser
	if err := repository.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// ListUsers 分页查询用户列表（不含密码）。
func ListUsers(page, pageSize int) ([]model.SysUser, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var total int64
	if err := repository.DB.Model(&model.SysUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.SysUser
	err := repository.DB.Model(&model.SysUser{}).
		Select("id", "user_name", "email", "address", "create_time", "update_time", "user_id").
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// 新建用户
func CreateUser(user interface{}) error { // 这里使用 interface{} 是因为 user 可以是 model.SysUser 或者 model.RegisterUser
	// GORM Create 需要传入指针类型
	if err := repository.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// 更新用户
func UpdateUser(user *model.SysUser) error {
	if err := repository.DB.Model(&model.SysUser{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// 删除用户
func DeleteUser(id uint) error {
	if err := repository.DB.Delete(&model.SysUser{}, id).Error; err != nil {
		return err
	}
	return nil
}
