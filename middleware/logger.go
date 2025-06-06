package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// Logger 创建日志中间件
// 返回 gin.HandlerFunc 用于集成到Gin框架
// 功能特性：
//   - 基于logrus的日志记录
//   - 日志文件轮转（按天分割）
//   - 多级别日志分类记录
//   - 请求元数据采集（响应时间、状态码、客户端信息等）
//
// 建议改进：
//   - 可扩展为JSON格式日志
//   - 增加日志压缩功能
//   - 支持动态配置日志级别
func Logger() gin.HandlerFunc {
	// 基础日志文件路径（需确保log目录存在）
	filePath := "log/log"
	linkName := "latest_log.log"

	// 尝试创建/打开基础日志文件（用于初始化）
	scr, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("日志文件初始化失败:", err)
	}

	// 创建logrus实例
	logger := logrus.New()
	logger.Out = scr                   // 设置默认输出
	logger.SetLevel(logrus.DebugLevel) // 设置日志记录级别（DEBUG及以上）

	// 配置日志轮转策略
	// 参数说明：
	//   - WithMaxAge: 日志保留时长（7天）
	//   - WithRotationTime: 日志分割间隔（24小时）
	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log", // 按日期格式分割日志
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(24*time.Hour),
		retalog.WithLinkName(linkName),
	)

	// 创建多级别日志写入映射
	// 所有级别日志均写入同一轮转文件
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	// 创建日志格式钩子
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 自定义时间格式
	})
	//logger.AddHook(Hook)

	// 创建控制台输出钩子（新增部分）
	consoleFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,  // 启用颜色
		FullTimestamp:   true,  // 显示完整时间
		DisableSorting:  false, // 保持字段顺序
		PadLevelText:    true,  // 对齐日志级别
	}

	consoleHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  os.Stdout,
		logrus.ErrorLevel: os.Stderr,
		logrus.WarnLevel:  os.Stdout,
		logrus.DebugLevel: os.Stdout,
		logrus.PanicLevel: os.Stderr,
	}, consoleFormatter)

	// 添加双钩子（文件和控制台）
	logger.AddHook(Hook)
	logger.AddHook(consoleHook)

	// 中间件处理函数
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := generateRequestID()

		// 创建带请求ID的上下文
		ctx := context.WithValue(c.Request.Context(), "RequestID", requestID)
		c.Request = c.Request.WithContext(ctx)
		// 记录请求开始时间
		startTime := time.Now()
		//// 新增：记录请求参数
		//var requestBody string
		//if c.Request.Body != nil {
		//	bodyBytes, _ := c.GetRawData()
		//	requestBody = string(bodyBytes)
		//	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		//}
		// 捕获请求体（支持重复读取）
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		// 处理请求（继续后续中间件）
		c.Next()

		// 计算请求处理耗时
		stopTime := time.Since(startTime).Milliseconds()
		spendTime := fmt.Sprintf("%d ms", stopTime)

		// 获取系统主机名
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}

		// 收集请求元数据
		statusCode := c.Writer.Status()    // HTTP状态码
		clientIp := c.ClientIP()           // 客户端IP
		userAgent := c.Request.UserAgent() // 用户代理信息
		dataSize := c.Writer.Size()        // 响应数据大小
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method   // 请求方法
		path := c.Request.RequestURI // 请求路径
		query := c.Request.URL.Query()
		body := sanitizeBody(string(bodyBytes)) // 脱敏处理敏
		// 构造日志条目
		entry := logger.WithFields(logrus.Fields{
			"RequestID": requestID,
			//"Latency":   time.Since(startTime).String(),//处理耗时
			"HostName":  hostName,   // 服务器标识
			"status":    statusCode, // HTTP状态码
			"SpendTime": spendTime,  // 处理耗时
			"Ip":        clientIp,   // 客户端IP
			"Method":    method,     // 请求方法
			"Path":      path,       // 请求路径
			"DataSize":  dataSize,   // 响应数据大小（字节）
			"Agent":     userAgent,  // 客户端信息
			"Query":     query,
			"Body":      body, // 脱敏处理敏
		})

		// 根据状态码分级记录
		switch {
		case len(c.Errors) > 0: // Gin框架错误
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		case statusCode >= 500: // 服务器端错误
			entry.Error()
		case statusCode >= 400: // 客户端错误
			entry.Warn()
		default: // 成功请求（1xx, 2xx, 3xx）
			entry.Info()
		}
	}
}

// 生成唯一请求ID（示例）
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 敏感字段脱敏（如密码）
func sanitizeBody(body string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err == nil {
		if _, ok := data["password"]; ok {
			data["password"] = "******"
		}
		sanitized, _ := json.Marshal(data)
		return string(sanitized)
	}
	return body
}
