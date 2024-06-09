package model

import (
	"gorm.io/gorm"
)

type User struct {
	ID    int32  `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Age   int32
	Email string `gorm:"size:100;unique"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
