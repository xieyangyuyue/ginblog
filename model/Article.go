package model

import (
	"context"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid;references:ID"`
	gorm.Model
	Title        string `gorm:"type:varchar(100);not null" json:"title"`
	Cid          int    `gorm:"type:int;not null" json:"cid"`
	Desc         string `gorm:"type:varchar(200)" json:"desc"`
	Content      string `gorm:"type:longtext" json:"content"`
	Img          string `gorm:"type:varchar(100)" json:"img"`
	CommentCount int    `gorm:"type:int;not null;default:0" json:"comment_count"`
	ReadCount    int    `gorm:"type:int;not null;default:0" json:"read_count"`
}

// CreateArt 新增文章
func CreateArt(ctx context.Context, data *Article) int {
	err := db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// GetCateArt 查询分类下的所有文章
func GetCateArt(ctx context.Context, id int, pageSize int, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64

	err := db.WithContext(ctx).Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"cid =?", id).Find(&cateArtList).Error
	db.WithContext(ctx).Model(&cateArtList).Where("cid =?", id).Count(&total)
	if err != nil {
		return nil, errmsg.ErrorCateNotExist, 0
	}
	return cateArtList, errmsg.Success, total
}

// GetArtInfo 查询单个文章
func GetArtInfo(ctx context.Context, id int) (Article, int) {
	var art Article
	err := db.WithContext(ctx).Where("id = ?", id).Preload("Category").First(&art).Error
	db.WithContext(ctx).Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if err != nil {
		return art, errmsg.ErrorCateNotExist
	}
	return art, errmsg.Success
}

// GetArt 查询文章列表
func GetArt(ctx context.Context, pageSize int, pageNum int) ([]Article, int, int64) {
	//var articleList []Article
	//var err error
	//var total int64
	//
	//err = db.Select("article.id, title, img, created_at, updated_at, `desc`, comment_count, read_count, category.name").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Joins("Category").Find(&articleList).Error
	//// 单独计数
	//db.Model(&articleList).Count(&total)
	//if err != nil {
	//	return nil, errmsg.Error, 0
	//}
	//return articleList, errmsg.Success, total

	var cateArtList []Article
	var total int64

	err := db.WithContext(ctx).Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cateArtList).Count(&total).Error
	if err != nil {
		return nil, errmsg.Error, 0
	}
	return cateArtList, errmsg.Success, total

}

// SearchArticle 搜索文章标题
func SearchArticle(ctx context.Context, title string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var err error
	var total int64
	err = db.WithContext(ctx).Select("article.id,title, img, created_at, updated_at, `desc`, comment_count, read_count, Category.name").Order("Created_At DESC").Joins("Category").Where("title LIKE ?",
		title+"%",
	).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&articleList).Error
	//单独计数
	db.WithContext(ctx).Model(&articleList).Where("title LIKE ?",
		title+"%",
	).Count(&total)

	if err != nil {
		return nil, errmsg.Error, 0
	}
	return articleList, errmsg.Success, total
}

// EditArt 编辑文章
func EditArt(ctx context.Context, id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img

	err := db.WithContext(ctx).Model(&art).Where("id = ? ", id).Updates(&maps).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// DeleteArt 删除文章
func DeleteArt(ctx context.Context, id int) int {
	var art Article
	err := db.WithContext(ctx).Where("id = ? ", id).Delete(&art).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}
