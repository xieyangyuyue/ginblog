package utils

import (
	"fmt"
	"gopkg.in/ini.v1" // 用于读取INI格式的配置文件
)

// 全局配置变量（包级作用域）
var (
	// AppMode 服务器配置
	AppMode  string // 应用模式（debug/release）
	HttpPort string // HTTP服务端口
	JwtKey   string // JWT令牌加密密钥

	// DbHost 数据库配置
	DbHost     string // 数据库主机地址
	DbPort     string // 数据库端口
	DbUser     string // 数据库用户名
	DbPassWord string // 数据库密码（注意：建议从环境变量获取敏感信息）
	DbName     string // 数据库名称

	// Zone 七牛云存储配置
	Zone       int    // 存储区域编号（1:华东 2:华北 3:华南）
	AccessKey  string // 七牛云AccessKey
	SecretKey  string // 七牛云SecretKey
	Bucket     string // 存储空间名称
	QiniuSever string // 七牛云服务地址
)

// 包初始化函数（自动执行）
func init() {
	// 加载配置文件（路径：config/config.ini）
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	// 分别加载不同配置模块
	LoadServer(file) // 加载服务器配置
	LoadData(file)   // 加载数据库配置
	LoadQiniu(file)  // 加载七牛云配置
}

// LoadServer 加载服务器配置模块
func LoadServer(file *ini.File) {
	// 从[server]区块读取配置项，MustString设置默认值
	section := file.Section("server")
	AppMode = section.Key("AppMode").MustString("debug")    // 默认开发模式
	HttpPort = section.Key("HttpPort").MustString(":3000")  // 默认端口3000
	JwtKey = section.Key("JwtKey").MustString("89js82js72") // 默认测试用密钥（生产环境必须修改！）
}

// LoadData 加载数据库配置模块
func LoadData(file *ini.File) {
	section := file.Section("database")
	DbHost = section.Key("DbHost").MustString("localhost") // 默认本地数据库
	DbPort = section.Key("DbPort").MustString("3306")      // 默认MySQL端口
	DbUser = section.Key("DbUser").MustString("ginblog")   // 默认用户名
	DbPassWord = section.Key("DbPassWord").String()        // 密码无默认值（必须配置）
	DbName = section.Key("DbName").MustString("ginblog")   // 默认数据库名
}

// LoadQiniu 加载七牛云配置模块
func LoadQiniu(file *ini.File) {
	section := file.Section("qiniu")
	Zone = section.Key("Zone").MustInt(1)           // 默认华东区域
	AccessKey = section.Key("AccessKey").String()   // AccessKey（必须配置）
	SecretKey = section.Key("SecretKey").String()   // SecretKey（必须配置）
	Bucket = section.Key("Bucket").String()         // 存储桶名称（必须配置）
	QiniuSever = section.Key("QiniuSever").String() // 服务地址（必须配置）
}
