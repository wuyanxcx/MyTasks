package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "newProject/task4/proto"
)

func main() {
	// 创建一个到grpc服务器的客户端连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 创建客户端接口
	grpcClient := pb.NewUserServiceClient(conn)
	// 创建一个默认的gin web服务器实例
	r := gin.Default()

	// 查
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// 创建请求
		req := &pb.GetUserInfoRequest{Id: int32(userId)}
		res, err := grpcClient.GetUserInfo(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	// 增
	r.POST("/user", func(c *gin.Context) {
		var req pb.CreateUserRequest
		// 将请求体中的json数据绑定到req变量
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := grpcClient.CreateUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	// 改
	r.PUT("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req pb.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Id = int32(userId)

		res, err := grpcClient.UpdateUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	// 删
	r.DELETE("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		req := &pb.DeleteUserRequest{Id: int32(userId)}
		res, err := grpcClient.DeleteUser(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.Run(":8080")
}
