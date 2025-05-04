package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"ginblog/utils/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddUser 添加用户
// @Summary 添加用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.User true "用户信息"
// @Success 200 {object} gin.H "{"status": 200, "message": "操作成功"}"
// @Failure 400 {object} gin.H "{"error": "无效的请求数据"}"
// @Router /api/v1/user [post]
// 处理 POST 请求，接收 JSON 格式用户数据，创建新用户
func AddUser(c *gin.Context) {
	// 获取请求上下文
	ctx := c.Request.Context()
	var data model.User
	var msg string
	var validCode int
	// 绑定请求中的 JSON 数据到 User 结构体
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	msg, validCode = validator.Validate(&data)
	if validCode != errmsg.Success {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		c.Abort()
		return
	}
	// 检查用户名是否已存在
	code := model.CheckUser(ctx, data.Username)
	if code == errmsg.Success {
		// 用户名未占用，执行创建操作
		model.CreateUser(ctx, &data)
	}
	if code == errmsg.ErrorUsernameUsed {
		code = errmsg.ErrorUsernameUsed
	}
	// 返回操作结果（状态码和消息）
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetUsers 查询用户列表
// @Summary 查询用户列表
// @Description 获取所有用户列表（分页待实现）
// @Tags 用户管理
// @Produce json
// @Success 200 {object} gin.H "{"status": 200, "data": []model.User}"
// @Router /api/v1/users [get]
func GetUsers(c *gin.Context) {
	// TODO: 分页查询用户列表逻辑
	ctx := c.Request.Context()
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetUsers(ctx, username, pageSize, pageNum)

	code := errmsg.Success
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// EditUser 编辑用户
// @Summary 编辑用户
// @Description 根据ID修改用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body model.User true "修改后的用户信息"
// @Success 200 {object} gin.H "{"status": 200, "message": "操作成功"}"
// @Failure 400 {object} gin.H "{"error": "无效的请求数据"}"
// @Router /api/v1/user/{id} [put]
func EditUser(c *gin.Context) {
	// TODO: 根据ID修改用户信息逻辑
	ctx := c.Request.Context()
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.CheckUpUser(ctx, id, data.Username)
	if code == errmsg.Success {
		model.EditUser(ctx, id, &data)
	}
	if code == errmsg.ErrorUsernameUsed {
		c.Abort()
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据ID软删除用户
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} gin.H "{"status": 200, "message": "操作成功"}"
// @Router /api/v1/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	// TODO: 根据ID软删除用户逻辑
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	code := model.DeleteUser(ctx, id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
