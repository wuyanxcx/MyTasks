package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"newProject/task4/modle"
	pb "newProject/task4/proto"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *gorm.DB
}

func (s *server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	var user model.User
	if err := s.db.First(&user, req.GetId()).Error; err != nil {
		return nil, err
	}
	return &pb.GetUserInfoResponse{
		Id:    user.ID,
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := model.User{
		Name:  req.GetName(),
		Age:   req.GetAge(),
		Email: req.GetEmail(),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{Id: user.ID}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var user model.User
	if err := s.db.First(&user, req.GetId()).Error; err != nil {
		return nil, err
	}
	user.Name = req.GetName()
	user.Age = req.GetAge()
	user.Email = req.GetEmail()
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{Success: true}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.db.Delete(&model.User{}, req.GetId()).Error; err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// 创建表 自动迁移（把结构体和数据表进行对应）
	model.AutoMigrate(db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建一个新的grpc服务器实例
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterUserServiceServer(s, &server{db: db})
	log.Printf("server listening at %v", lis.Addr())
	// 启动服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
