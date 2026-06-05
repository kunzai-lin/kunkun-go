package model

type RegisterUser struct {
	UserName     string `json:"username" binding:"required" gorm:"column:user_name;type:varchar(255)"`
	UserPassword string `json:"password" binding:"required" gorm:"column:user_password;type:varchar(255)"`
}

// TableName 指定 RegisterUser 结构体对应的数据库表名
func (*RegisterUser) TableName() string {
	return "sys_user"
}


