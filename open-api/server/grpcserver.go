package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-gateway/open-api/internal/api"
	"log"
	"net"
)

type Server_wc struct {
	api.UnimplementedWenchangchainServer
}
type Server_sz struct {
	api.UnimplementedShenzhouServer
}

func (*Server_sz) CreateNftClass(context.Context, *api.CreateNftRequest) (*api.CreateNftResponse, error) {
	fmt.Println("shenzhou CreateNftClass grpc successfully")
	response := &api.CreateNftResponse{
		Data: "shenzhou create 启动正常",
	}
	return response, nil
}

func (*Server_sz) ClassByID(context.Context, *api.ClassByIDRequest) (*api.ClassByIDResponse, error) {
	fmt.Println(" shenzhou grpc 处理方法调用 successfully")
	response := &api.ClassByIDResponse{
		Name: "shenzhou CreateByID 启动正常",
	}
	return response, nil
}

func (*Server_wc) CreateNftClass(context.Context, *api.CreateNftRequest) (*api.CreateNftResponse, error) {
	fmt.Println("wenchan CreateNftClass grpc 处理方法调用 successfully")
	response := &api.CreateNftResponse{
		Data: "wenchan create grpc 启动正常",
	}
	return response, nil
}
func (*Server_wc) ClassByID(context.Context, *api.ClassByIDRequest) (*api.ClassByIDResponse, error) {
	fmt.Println("wenchan ClassByID grpc 处理方法调用 successfully")
	response := &api.ClassByIDResponse{
		ID:   "12",
		Name: "shihengtest",
	}
	return response, nil
}
func main() {
	// 创建监听器
	lisWenchang, err := net.Listen("tcp", ":50051") //向外暴露grpc服务端口
	if err != nil {
		log.Fatalln(err)
	}
	lisShenzhou, err := net.Listen("tcp", ":50052") //向外暴露grpc服务端口
	if err != nil {
		log.Fatalln(err)
	}
	sWenchang := grpc.NewServer() // 创建服务器实例

	sShenzho := grpc.NewServer() // 创建服务器实例
	api.RegisterWenchangchainServer(sWenchang, &Server_wc{})

	api.RegisterShenzhouServer(sShenzho, &Server_sz{})
	// 启动grpc服务器
	go func() {
		if err := sWenchang.Serve(lisWenchang); err != nil {
			log.Println("wenchang server error:", err)
		}
	}()

	go func() {
		if err := sShenzho.Serve(lisShenzhou); err != nil {
			log.Println("shenzhou server error:", err)
		}
	}()

	// 无限循环，保持服务器运行
	for {
		select {}
	}
}
