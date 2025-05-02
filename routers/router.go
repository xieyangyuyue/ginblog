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
		//新增用户
		router.POST("user/add", v1.AddUser)
		//查询用户列表
		router.GET("users", v1.GetUsers)
		//编辑用户
		router.PUT("user/:id", v1.EditUser)
		//删除用户
		router.DELETE("user/:id", v1.DeleteUser)

		// 分类模块的路由接口
		//查询分类列表
		router.GET("category", v1.GetCate)
		//查询具体分类
		router.GET("category/:id", v1.GetCateInfo)
		//添加分类
		router.POST("category/add", v1.AddCategory)
		//编辑分类
		router.PUT("category/:id", v1.EditCate)
		//删除分类
		router.DELETE("category/:id", v1.DeleteCate)

		//文章模块的路由接口
		//查询文章列表
		router.GET("article", v1.GetArt)
		//查询单个文章信息
		router.GET("article/info/:id", v1.GetArtInfo)
		//查询分类下的所有文章
		router.GET("article/list/:id", v1.GetCateArt)
		//新增文章
		router.POST("article/add", v1.AddArticle)
		//编辑文章
		router.PUT("article/:id", v1.EditArt)
		//删除文章
		router.DELETE("article/:id", v1.DeleteArt)
	}

	// 启动HTTP服务（从配置中读取端口号）
	// utils.HttpPort 示例值：":8080"（冒号+端口号格式）
	err := r.Run(utils.HttpPort)
	if err != nil {
		return
	}
}
