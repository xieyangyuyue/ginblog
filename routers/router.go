package routers

import (
	v1 "ginblog/api/v1"
	"ginblog/utils"
	"github.com/gin-gonic/gin"
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
		//用户模块的路由接口
		router.POST("user/add", v1.AddUser)
		router.GET("users", v1.GetUsers)
		router.PUT("user/:id", v1.EditUser)
		router.DELETE("user/:id", v1.DeleteUser)

		// 分类模块的路由接口
		router.GET("category", v1.GetCate)
		router.GET("category/:id", v1.GetCateInfo)
		router.POST("category/add", v1.AddCategory)
		router.PUT("category/:id", v1.EditCate)
		router.DELETE("category/:id", v1.DeleteCate)

		//文章模块的路由接口
		router.GET("article", v1.GetArt)
		router.POST("article/add", v1.AddArticle)
		router.PUT("article/:id", v1.EditArt)
		router.DELETE("article/:id", v1.DeleteArt)
	}

	// 启动HTTP服务（从配置中读取端口号）
	// utils.HttpPort 示例值：":8080"（冒号+端口号格式）
	err := r.Run(utils.HttpPort)
	if err != nil {
		return
	}
}
