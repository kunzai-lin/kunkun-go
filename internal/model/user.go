package model

import "time"

type SysUser struct {
	ID           uint      `gorm:"primaryKey;column:id;type:bigint(20);"`
	UserName     string    `gorm:"column:user_name;type:varchar(255);"`
	UserPassword string    `gorm:"column:user_password;type:varchar(255);"`
	Email        string    `gorm:"column:email;type:varchar(255);"`
	Address      string    `gorm:"column:address;type:varchar(255);"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime"`
	UserID       uint      `gorm:"column:user_id;type:bigint(20);"`
}

// TableName 指定 SysUser 结构体对应的数据库表名
func (*SysUser) TableName() string {
	return "sys_user"
}
