package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "newProject/task2/proto"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// 模拟用户数据
	user := &pb.GetUserInfoResponse{
		Id:    req.GetId(),
		Name:  "John Doe",
		Age:   30,
		Email: "john.doe@example.com",
	}
	return user, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
