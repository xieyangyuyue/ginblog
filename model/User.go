package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

// User 用户模型（对应数据库表）
type User struct {
	gorm.Model        // 内嵌 gorm.Model，包含字段 ID、CreatedAt、UpdatedAt、DeletedAt
	Username   string `gorm:"type:varchar(20);not null"`                  // 用户名，数据库约束：长度20，非空
	Password   string `gorm:"type:varchar(500);not null" json:"password"` // 密码，存储加密后的值，非空
	Role       int    `gorm:"type:int;DEFAULT:2" json:"role"`             // 角色，1-管理员，2-普通用户，默认值2
}

// CheckUser 检查用户名是否存在
// 参数 name: 待检查的用户名
// 返回值 code: 错误码（成功或错误类型）
func CheckUser(name string) (code int) {
	var user User
	// 查询数据库中是否存在同名用户（只查询ID字段）
	db.Select("id").Where("username = ?", name).First(&user)
	if user.ID > 0 {
		return errmsg.ErrorUsernameUsed // 返回错误码 1001（此处逻辑可能需要调整，通常应为“用户名已存在”）
	}
	return errmsg.Success // 返回成功码 200
}

// CreateUser 创建新用户
// 参数 data: 用户数据指针
// 返回值 int: 错误码（成功或错误类型）
func CreateUser(data *User) int {
	// 密码加密逻辑（示例中暂时被注释）
	// data.Password = ScryptPw(data.Password)

	// 执行数据库插入操作
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.Error // 返回错误码 500
	}
	return errmsg.Success // 返回成功码 200
}

// GetUsers 查询用户列表
func GetUsers(username string, pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64

	if username != "" {
		db.Select("id,username,role,created_at").Where(
			"username LIKE ?", username+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		db.Model(&users).Where(
			"username LIKE ?", username+"%",
		).Count(&total)
		return users, total
	}
	err := db.Select("id,username,role,created_at").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	db.Model(&users).Count(&total)

	if err != nil {
		return users, 0
	}
	return users, total
}
