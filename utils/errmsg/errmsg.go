// Package errmsg 定义应用程序的错误码和错误消息映射
package errmsg

// 应用级通用状态码
const (
	Success = 200 // 成功状态码
	Error   = 500 // 通用错误状态码
)

// 用户模块错误码 (1001-1008)
const (
	ErrorUsernameUsed   = 1001 + iota // 用户名已被使用
	ErrorPasswordWrong                // 密码不正确
	ErrorUserNotExist                 // 用户不存在
	ErrorTokenExist                   // TOKEN不存在
	ErrorTokenRuntime                 // TOKEN已过期
	ErrorTokenWrong                   // TOKEN无效
	ErrorTokenTypeWrong               // TOKEN类型错误
	ErrorUserNoRight                  // 用户无权限
)

// ErrorArtNotExist 文章模块错误码 (2001)
const (
	ErrorArtNotExist = 2001 // 文章不存在
)

// 分类模块错误码 (3001-3002)
const (
	ErrorCatenameUsed = 3001 + iota // 分类名称已存在
	ErrorCateNotExist               // 分类不存在
)

// codeMsg 错误码与错误信息的映射表
var codeMsg = map[int]string{
	Success: "OK",
	Error:   "内部错误",

	// 用户模块
	ErrorUsernameUsed:   "用户名已被占用",
	ErrorPasswordWrong:  "密码验证失败",
	ErrorUserNotExist:   "用户不存在",
	ErrorTokenExist:     "身份令牌缺失，请重新登录",
	ErrorTokenRuntime:   "身份令牌已过期，请重新登录",
	ErrorTokenWrong:     "无效的身份令牌",
	ErrorTokenTypeWrong: "非法的令牌格式",
	ErrorUserNoRight:    "用户权限不足",

	// 文章模块
	ErrorArtNotExist: "指定文章不存在",

	// 分类模块
	ErrorCatenameUsed: "分类名称已存在",
	ErrorCateNotExist: "指定分类不存在",
}

// GetErrMsg 根据错误码获取对应的错误信息
// 当code不存在时返回空字符串
func GetErrMsg(code int) string {
	return codeMsg[code]
}

// GetErrMsgWithDefault 根据错误码获取错误信息，不存在时返回默认消息
func GetErrMsgWithDefault(code int, defaultMsg string) string {
	if msg, ok := codeMsg[code]; ok {
		return msg
	}
	return defaultMsg
}
