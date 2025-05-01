package routers

import (
	"ginblog/utils"       // 自定义配置工具包
	"github.com/gin-gonic/gin"  // Gin Web框架
	"net/http"            // HTTP协议支持
)

// InitRouter 初始化路由并启动HTTP服务
// 返回值：无
func InitRouter() {
	// 设置Gin运行模式（从配置中读取）
	// utils.AppMode 可能的值：debug/test/release
	gin.SetMode(utils.AppMode) 

	// 创建默认路由引擎（自带Logger和Recovery中间件）
	r := gin.Default()

	// 创建API路由分组（版本控制）
	// 所有路由将以 /api/v1/ 作为前缀
	router := r.Group("api/v1")
	{
		// 示例测试路由（GET请求）
		// 访问路径：/api/v1/hello
		router.GET("hello", func(c *gin.Context) {
			// 返回JSON格式响应
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",  // 示例响应内容
			})
		})

		// 后续可以在此处添加更多路由：
		// router.POST("/articles", 创建文章)
		// router.GET("/articles", 获取文章列表)
		// router.PUT("/articles/:id", 更新文章)
	}

	// 启动HTTP服务（从配置中读取端口号）
	// utils.HttpPort 示例值：":8080"（冒号+端口号格式）
	r.Run(utils.HttpPort) 
}