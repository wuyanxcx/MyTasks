package models

import (
	"gorm.io/gorm"
)

// User 定义用户模型
type User struct {
	gorm.Model
	Name  string
	Email string
}
