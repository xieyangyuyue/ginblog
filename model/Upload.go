// Package model Upload.go
package model

import (
	"context"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
)

// 七牛云配置全局变量
var (
	Zone      = utils.Zone       // 存储区域编号（1-华东，2-华北，3-华南）
	AccessKey = utils.AccessKey  // 七牛云AccessKey
	SecretKey = utils.SecretKey  // 七牛云SecretKey
	Bucket    = utils.Bucket     // 存储空间名称
	ImgUrl    = utils.QiniuSever // 七牛云图片访问域名
)

// UpLoadFile 上传文件到七牛云存储
// 参数：
//
//	file - 要上传的文件对象
//	fileSize - 文件大小（字节）
//
// 返回值：
//
//	string - 文件访问URL
//	int - 状态错误码
func UpLoadFile(file multipart.File, fileSize int64) (string, int) {
	// 创建上传策略
	putPolicy := storage.PutPolicy{
		Scope: Bucket, // 指定目标存储空间
	}
	// 创建Mac对象用于签名
	mac := qbox.NewMac(AccessKey, SecretKey)
	// 生成上传凭证
	upToken := putPolicy.UploadToken(mac)

	// 获取存储配置
	cfg := setConfig()

	// 创建表单上传对象
	formUploader := storage.NewFormUploader(&cfg)
	// 初始化返回结构
	ret := storage.PutRet{}
	// 不需要额外参数
	putExtra := storage.PutExtra{}

	// 执行上传操作（不指定存储文件名）
	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", errmsg.Error // 返回上传错误
	}
	// 拼接完整访问URL
	url := ImgUrl + ret.Key
	return url, errmsg.Success
}

// setConfig 配置七牛云存储区域
// 返回值：storage.Config - 存储配置对象
func setConfig() storage.Config {
	return storage.Config{
		Zone:          selectZone(Zone), // 根据配置选择存储区域
		UseCdnDomains: false,            // 不使用CDN域名
		UseHTTPS:      false,            // 不使用HTTPS
	}
}

// selectZone 根据区域编号选择七牛云存储区域
// 参数：id - 区域编号（1-华东，2-华北，3-华南）
// 返回值：*storage.Zone - 对应的存储区域指针
func selectZone(id int) *storage.Zone {
	switch id {
	case 1:
		return &storage.ZoneHuadong // 华东区域
	case 2:
		return &storage.ZoneHuabei // 华北区域
	case 3:
		return &storage.ZoneHuanan // 华南区域
	default:
		return &storage.ZoneHuadong // 默认华东区域
	}
}
