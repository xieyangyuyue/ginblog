package routers

import (
	v1 "ginblog/api/v1"
	"ginblog/middleware"
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
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口
		//新增用户
		auth.POST("user/add", v1.AddUser)
		//编辑用户
		auth.PUT("user/:id", v1.EditUser)
		//删除用户
		auth.DELETE("user/:id", v1.DeleteUser)

		//修改密码
		//auth.PUT("admin/changepw/:id", v1.ChangeUserPassword)

		// 分类模块的路由接口
		//添加分类
		auth.POST("category/add", v1.AddCategory)
		//编辑分类
		auth.PUT("category/:id", v1.EditCate)
		//删除分类
		auth.DELETE("category/:id", v1.DeleteCate)

		// 文章模块的路由接口
		//新增文章
		auth.POST("article/add", v1.AddArticle)
		//编辑文章
		auth.PUT("article/:id", v1.EditArt)
		//删除文章
		auth.DELETE("article/:id", v1.DeleteArt)
		// 上传文件
		auth.POST("upload", v1.UpLoad)
		//// 更新个人设置
		//auth.GET("admin/profile/:id", v1.GetProfile)
		//auth.PUT("profile/:id", v1.UpdateProfile)
		//// 评论模块
		//auth.GET("comment/list", v1.GetCommentList)
		//auth.DELETE("delcomment/:id", v1.DeleteComment)
		//auth.PUT("checkcomment/:id", v1.CheckComment)
		//auth.PUT("uncheckcomment/:id", v1.UncheckComment)
	}

	router := r.Group("api/v1")
	{
		//用户模块的路由接口
		//查询用户列表
		router.GET("users", v1.GetUsers)

		// 分类模块的路由接口
		//查询分类列表
		router.GET("category", v1.GetCate)
		//查询具体分类
		router.GET("category/:id", v1.GetCateInfo)

		//文章模块的路由接口
		//查询文章列表
		router.GET("article", v1.GetArt)
		//查询单个文章信息
		router.GET("article/info/:id", v1.GetArtInfo)
		//查询分类下的所有文章
		router.GET("article/list/:id", v1.GetCateArt)

		// 登录控制模块
		router.POST("login", v1.Login)
		router.POST("loginfront", v1.LoginFront)

	}

	// 启动HTTP服务（从配置中读取端口号）
	// utils.HttpPort 示例值：":8080"（冒号+端口号格式）
	err := r.Run(utils.HttpPort)
	if err != nil {
		return
	}
}
