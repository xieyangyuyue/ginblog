// Package validator  数据验证模块包
package validator

import (
	"fmt"
	"ginblog/utils/errmsg"                                           // 自定义错误码
	"github.com/go-playground/locales/zh_Hans_CN"                    // 中文本地化包
	unTrans "github.com/go-playground/universal-translator"          // 通用翻译器
	"github.com/go-playground/validator/v10"                         // 主验证库
	zhTrans "github.com/go-playground/validator/v10/translations/zh" // 中文翻译
	"reflect"                                                        // 反射包用于获取结构体标签
)

// Validate 数据验证函数
// 参数 data: 需要验证的结构体数据
// 返回值: (错误信息, 错误码)
func Validate(data any) (string, int) {
	// 创建新的验证器实例
	validate := validator.New()

	// 创建中文翻译器
	uni := unTrans.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN") // 获取中文翻译器

	// 注册默认中文翻译
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("翻译器注册错误:", err)
	}

	// 注册标签名函数：使用结构体的 "label" 标签作为字段名称
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label") // 获取结构体字段的 label 标签值
		return label
	})

	// 执行结构体验证
	err = validate.Struct(data)
	if err != nil {
		// 类型断言转换为验证错误集合
		for _, v := range err.(validator.ValidationErrors) {
			// 返回第一个错误的中文翻译和错误码
			return v.Translate(trans), errmsg.Error
		}
	}

	// 验证成功返回空字符串和成功码
	return "", errmsg.Success
}
