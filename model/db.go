package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	"ginblog/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDb() {
	// 构建DSN（数据源名称）
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.DbUser,
		utils.DbPassWord,
		utils.DbHost,
		utils.DbPort,
		utils.DbName,
	)
	// 配置GORM日志（更新为以下内容）
	newLogger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info, // 显示所有SQL日志
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true, // 显示完整参数（非占位符形式）
			Colorful:                  true, // 控制台彩色输出
		},
	)
	// 初始化GORM配置
	config := &gorm.Config{
		Logger: newLogger,
		//Logger: logger.Default.LogMode(logger.Silent), // 静默日志模式
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表名
		},
		SkipDefaultTransaction:                   true,  // 禁用默认事务
		DisableForeignKeyConstraintWhenMigrating: false, // 注意这里保持外键约束
	}

	// 连接数据库
	var err error
	if db, err = gorm.Open(mysql.Open(dns), config); err != nil {
		log.Fatal("数据库连接失败: ", err)
		os.Exit(1)
	}
	// ✅ 正确位置：在数据库连接成功后执行测试查询
	if err := db.Debug().Exec("SELECT 1 + 1").Error; err != nil {
		log.Fatal("数据库连接测试失败: ", err)
		os.Exit(1)
	}
	// 获取底层SQL DB对象以设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接池失败: ", err)
		os.Exit(1)
	}

	// 连接池设置
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 建议适当延长生命周期（原示例的10秒太短）
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 新增空闲超时设置

	// 自动迁移（建议生产环境谨慎使用）
	// , &Profile{}, &Comment{}
	if err := db.AutoMigrate(&User{}, &Article{}, &Category{}); err != nil {
		log.Fatal("数据库迁移失败: ", err)
		os.Exit(1)
	}
}
