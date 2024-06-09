package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"newProject/task5/handler"
	"newProject/task5/model"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	model.AutoMigrate(db)

	r := gin.Default()
	announcementHandler := handler.NewAnnouncementHandler(db)

	r.POST("/announcements", announcementHandler.CreateAnnouncement)
	r.GET("/announcements", announcementHandler.GetAnnouncements)
	r.GET("/announcements/:id", announcementHandler.GetAnnouncement)
	r.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)
	r.DELETE("/announcements/:id", announcementHandler.DeleteAnnouncement)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
