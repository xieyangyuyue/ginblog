package middleware

import (
	"errors"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // 使用 JWT v5 版本
	"net/http"
	"strings"
)

// JWT 结构体，包含 JWT 签名密钥
type JWT struct {
	JwtKey []byte
}

// NewJWT 创建 JWT 实例，从 utils 包获取密钥
func NewJWT() *JWT {
	return &JWT{
		JwtKey: []byte(utils.JwtKey),
	}
}

// MyClaims 自定义 Claims 结构体，包含用户名和标准 Claims
type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 定义全局错误变量
var (
	TokenExpired     = errors.New("token已过期,请重新登录")
	TokenNotValidYet = errors.New("token尚未生效,请重新登录")
	TokenMalformed   = errors.New("token格式错误,请重新登录")
	TokenInvalid     = errors.New("无效的token,请重新登录")
)

// CreateToken 生成 JWT Token
// claims: 包含用户信息的自定义 Claims
// 返回值: 生成的 Token 字符串或错误
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

// ParseToken 解析并验证 Token，返回自定义 Claims
// tokenString: 待解析的 Token 字符串
// 返回值: 解析后的 Claims 或错误
func (j *JWT) ParseToken(tokenString string) (*MyClaims, error) {
	// 解析 Token 并验证签名
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})

	if err != nil {
		// JWT v5 错误处理：使用预定义错误类型
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}

	// 验证 Token 有效性并提取 Claims
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// JwtToken JWT 中间件，用于 Gin 路由的 Token 验证
// 返回值: Gin 中间件 HandlerFunc
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		tokenHeader := c.GetHeader("Authorization")

		if tokenHeader == "" {
			// Token 缺失错误
			code = errmsg.ErrorTokenExist
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 检查 Token 格式是否为 "Bearer <token>"
		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ErrorTokenTypeWrong
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 解析 Token
		j := NewJWT()
		claims, err := j.ParseToken(checkToken[1])
		if err != nil {
			// 根据错误类型映射错误码
			switch {
			case errors.Is(err, TokenExpired):
				code = errmsg.ErrorTokenRuntime
			case errors.Is(err, TokenNotValidYet), errors.Is(err, TokenMalformed), errors.Is(err, TokenInvalid):
				code = errmsg.ErrorTokenWrong
			default:
				code = errmsg.Error
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户名存入 Gin 上下文，供后续处理使用
		c.Set("username", claims.Username)
		c.Next()
	}
}
