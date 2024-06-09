package repository

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"newProject/task3/models"
	"testing"
)

func TestUserCRUD(t *testing.T) {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	//创建表 自动迁移（把结构体和数据表进行对应）
	db.AutoMigrate(&models.User{})

	// 创建用户
	user := models.User{Name: "John Doe", Email: "john@example.com"}
	if err := CreateUser(db, &user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 查询用户
	retrievedUser, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("retrieved user name doesn't match, got: %s, want: %s", retrievedUser.Name, user.Name)
	}

	// 更新用户
	user.Name = "Jane Doe"
	if err := UpdateUser(db, &user); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}
	updatedUser, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Fatalf("failed to get updated user: %v", err)
	}
	if updatedUser.Name != "Jane Doe" {
		t.Errorf("updated user name doesn't match, got: %s, want: %s", updatedUser.Name, "Jane Doe")
	}

	// 删除用户
	if err := DeleteUser(db, user.ID); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	_, err = GetUserByID(db, user.ID)

	if err == nil {
		t.Error("expected error when getting deleted user, got nil")
	}

	// 断开数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database connection", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatal("failed to close database connection", err)
	}
}

// CreateUser 创建一个新用户
func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

// GetUserByID 通过 ID 获取用户
func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(db *gorm.DB, user *models.User) error {
	return db.Model(user).Updates(user).Error
}

// DeleteUser 删除用户
func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&models.User{}, id).Error
}
