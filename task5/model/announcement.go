package model

import (
	"gorm.io/gorm"
	"time"
)

type Announcement struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Announcement{})
}
