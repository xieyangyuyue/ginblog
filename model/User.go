package model

import (
	"crypto/rand"
	"encoding/base64"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/scrypt"
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

// CheckUpUser 更新查询
func CheckUpUser(id int, name string) (code int) {
	var user User
	db.Select("id, username").Where("username = ?", name).First(&user)
	if user.ID == uint(id) {
		return errmsg.Success
	}
	if user.ID > 0 {
		return errmsg.ErrorUsernameUsed //1001
	}
	return errmsg.Success
}

// CreateUser 创建新用户
// 参数 data: 用户数据指针
// 返回值 int: 错误码（成功或错误类型）
func CreateUser(data *User) int {
	// 密码加密逻辑（示例中暂时被注释）
	//data.Password = ScryptPw(data.Password)

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

// EditUser 编辑用户信息
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err := db.Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// DeleteUser 删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id = ? ", id).Delete(&user).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// BeforeCreate 密码加密&权限控制
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	u.Role = 2
	return nil
}

// ScryptPw 使用scrypt算法安全地处理密码存储
// 参数:	password - 用户输入的明文密码字符串
// 返回值:string - 包含随机盐值和哈希值的Base64组合字符串 error  - 执行过程中的错误信息
func ScryptPw(password string) string {
	// 算法参数配置（符合OWASP推荐基准）
	const (
		N       = 32768 // CPU/Memory开销参数，每轮迭代次数（需权衡安全性与性能）
		r       = 8     // 内存块大小参数，每个块的大小（字节数）
		p       = 1     // 并行度参数，建议保持为1避免侧信道攻击
		KeyLen  = 32    // 输出密钥长度（32字节=256位，满足AES-256安全要求）
		saltLen = 16    // 盐值长度（16字节=128位，推荐最小值）
	)
	// 生成密码学安全的随机盐值
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "" // 封装底层错误信息
	}
	// 执行scrypt密钥派生运算
	hash, err := scrypt.Key(
		[]byte(password), // 明文密码转换为字节切片
		salt,             // 使用新生成的随机盐值
		N, r, p,          // 调优后的算法参数
		KeyLen, // 指定输出密钥长度
	)
	if err != nil {
		return ""
	}
	// 组合盐值与哈希值（存储时需要完整保留）
	combined := append(salt, hash...) // 前16字节为盐，后32字节为哈希
	// 使用URL安全的Base64编码（避免+/字符，适合数据库存储）
	return base64.URLEncoding.EncodeToString(combined)
}
