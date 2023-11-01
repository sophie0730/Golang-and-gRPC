package main

import (
	"context"
	"fmt"
	pb "grpc/hello-server/proto"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)


type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// return &pb.HelloResponse{ResponseMsg: "hello" + req.RequestName}, nil\
	fmt.Println("gogogo")
	return &pb.HelloResponse{ResponseMsg: "GO"}, nil
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	//開啟端口
	listen, _ := net.Listen("tcp" ,":9090")
    //創建grpc服務
	grpcServer := grpc.NewServer()
    //在grpc服務端去註冊我們自己編寫的服務
	pb.RegisterSayHelloServer(grpcServer, &server{})

    //啟動服務
	go func ()  {
		err:= grpcServer.Serve(listen)
		if  err != nil {
			fmt.Printf("failed to serve: %v", err)
			return
		}
	} ()


	r := gin.Default();
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H {
			"response": "hi",
		})
	})
	r.POST("/method2_response", func (c *gin.Context)  {
		c.JSON(200, gin.H {
			"response": "method 2: hello",
		})
	})

	r.Run(":3000");
}