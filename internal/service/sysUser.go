package service

import (
	"kunkun-go/internal/model"
	"kunkun-go/internal/repository"
)

// 获取用户信息
func GetUserInfo(id uint) (model.SysUser, error) {
	var user model.SysUser
	if err := repository.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
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
