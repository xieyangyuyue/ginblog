package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " `
	Password string `gorm:"type:varchar(500);not null" json:"password"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" `
}