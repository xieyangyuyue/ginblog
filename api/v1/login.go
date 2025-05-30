package v1

import (
	"ginblog/middleware"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// Login 后台登陆
func Login(c *gin.Context) {
	// 获取请求上下文
	ctx := c.Request.Context()
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var token string
	var code int

	formData, code = model.CheckLogin(ctx, formData.Username, formData.Password)

	if code == errmsg.Success {
		setToken(c, formData)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    formData.Username,
			"id":      formData.ID,
			"message": errmsg.GetErrMsg(code),
			"token":   token,
		})
	}

}

// LoginFront 前台登录
func LoginFront(c *gin.Context) {
	// 获取请求上下文
	ctx := c.Request.Context()
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var code int

	formData, code = model.CheckLoginFront(ctx, formData.Username, formData.Password)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    formData.Username,
		"id":      formData.ID,
		"message": errmsg.GetErrMsg(code),
	})
}

// token生成函数
func setToken(c *gin.Context, user model.User) {
	j := middleware.NewJWT()
	claims := middleware.MyClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "GinBlog",
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.Error,
			"message": errmsg.GetErrMsg(errmsg.Error),
			"token":   token,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    user.Username,
		"id":      user.ID,
		"message": errmsg.GetErrMsg(200),
		"token":   token,
	})
	return
}
