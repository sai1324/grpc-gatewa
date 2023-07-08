package main

//服务器文件
import (
	"context"
	"google.golang.org/grpc"
	"grpc-gateway/test/proto/server"
	"log"
	"net"
)

type Server struct {
	server.UnimplementedCreateNftClassServer
}

// 实现处理方法
func (*Server) CreateNftClass(context.Context, *server.CreateNftRequest) (*server.CreateNftResponse, error) {
	println("CreateNftClass grpc 处理程序成功调用")
	response := &server.CreateNftResponse{}
	return response, nil
}
func (*Server) ClassByID(context.Context, *server.ClassByIDRequest) (*server.ClassByIDResponse, error) {
	println("ClassById grpc 处理程序成功调用")
	response := &server.ClassByIDResponse{}
	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50054") //创建一个tcp监听器
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()                             //创建grpc服务器实例，用于注册和处理grpc服务的请求
	server.RegisterCreateNftClassServer(s, &Server{}) //将处理程序和服务器程序绑定
	if err := s.Serve(lis); err != nil {              //开始监听并处理来自grpc的请求
		log.Fatalln(err)
	}
	for {
		select {}
	}
}

//func (s *Server) CreateNftClass(ctx context.Context, request *internal.CreateNftRequest) (*internal.CreateNftResponse, error) {
//	println("CreateNftClass grpc 处理程序成功调用")
//	response := &internal.CreateNftResponse{}
//	return response, nil
//}
//
//func (s *Server) ClassByID(ctx context.Context, request *internal.ClassByIDRequest) (*internal.ClassByIDResponse, error) {
//	println("ClassById grpc 处理程序成功调用")
//	response := &internal.ClassByIDResponse{}
//	return response, nil
//}
//
//func (s *Server) mustEmbedUnimplementedCreateNftClassServer() {
//	println("1")
//}
//
//func (s *Server) print(ctx context.Context) {
//	fmt.Println("处理create nft grpc 请求逻辑")
//
//}
