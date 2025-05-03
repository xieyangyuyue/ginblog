// Package v1 upload.go
package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpLoad 图片上传API接口
// 路由处理函数，接收客户端上传的文件并返回存储结果
// c - Gin上下文对象，包含请求和响应信息
func UpLoad(c *gin.Context) {
	// 从表单中获取文件对象
	file, fileHeader, _ := c.Request.FormFile("file")
	// 获取文件大小
	fileSize := fileHeader.Size

	// 调用模型层上传文件
	url, code := model.UpLoadFile(file, fileSize)

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"status":  code,                   // 状态码
		"message": errmsg.GetErrMsg(code), // 状态消息
		"url":     url,                    // 文件访问URL
	})
}
