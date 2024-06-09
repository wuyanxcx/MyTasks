package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "newProject/task2/proto"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewUserServiceClient(conn)

	// 创建 Gin 服务器
	r := gin.Default()

	// 定义 HTTP 接口
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 转换 id 为 int32
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// 调用 gRPC 服务
		req := &pb.GetUserInfoRequest{Id: int32(userId)}
		res, err := grpcClient.GetUserInfo(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回用户信息
		c.JSON(http.StatusOK, res)
	})

	// 启动 Gin 服务器
	r.Run(":8080")
}
