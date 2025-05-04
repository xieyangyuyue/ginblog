package model

import (
	"context"
	"errors"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// CheckCategory 查询分类是否存在
// @name 传过来的name字符串
func CheckCategory(ctx context.Context, name string) (code int) {
	var cate Category
	db.WithContext(ctx).Select("id").Where("name = ?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ErrorCatenameUsed
	}
	return errmsg.Success
}

// CheckUpCategory 更新查询
func CheckUpCategory(ctx context.Context, id int, name string) (code int) {
	var cate Category
	db.WithContext(ctx).Select("id, name").Where("name = ?", name).First(&cate)
	if cate.ID == uint(id) {
		return errmsg.Success
	}
	if cate.ID > 0 {
		return errmsg.ErrorCatenameUsed //1001
	}
	return errmsg.Success
}

// CreateCate 新增分类
func CreateCate(ctx context.Context, data *Category) int {
	err := db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return errmsg.Error // 500
	}
	return errmsg.Success
}

// GetCateInfo 查询单个分类信息
func GetCateInfo(ctx context.Context, id int) (Category, int) {
	var cate Category
	db.WithContext(ctx).Where("id = ?", id).First(&cate)
	return cate, errmsg.Success
}

// GetCate 查询分类列表
func GetCate(ctx context.Context, pageSize int, pageNum int) ([]Category, int64) {
	var cate []Category
	var total int64
	err := db.WithContext(ctx).Find(&cate).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.WithContext(ctx).Model(&cate).Count(&total)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0
	}
	return cate, total
}

// EditCate 编辑分类信息
func EditCate(ctx context.Context, id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	err := db.WithContext(ctx).Model(&cate).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// DeleteCate 删除分类
func DeleteCate(ctx context.Context, id int) int {
	var cate Category
	err := db.WithContext(ctx).Where("id = ? ", id).Delete(&cate).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}
