package main

import (
	"flag"
	pb "grpc-gateway/test/proto/server"
)

// 一个链服务写一个http的服务器文件
// grpc连接创建问题
// 每增加一个链服务都要重新写一个proto文件这是一定的，生成的网关文件（gw.internal.go) 中就生成了连接
// 网关中也会自己调用对应grpc处理方法的 我只需要写http的请求处理
import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	// 问题：如何从配置项中动态获取这个endpoint
	// 连接由谁去创建？  由网关自动创建，需要使用特殊配置可以使用option 只需要传入地址即可
	// 从配置项读出变成字符串----创建谁的连接是如何判断的 --用户传入的apikey决定
	grpcServerEndpoint = flag.String("grpc_server_endpoint", "localhost:50054", "") //grpc server endpoint
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

// 简单记录一下请求信息作为中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理中间件逻辑，例如记录请求信息
		log.Println("Received request:", r.Method, r.URL.Path)
		fmt.Println("Received request")
		fmt.Println("调用http前调用中间件")
		// 调用下一个处理程序 --转发到 mux
		next.ServeHTTP(w, r) //问题？如何知道下一个是转到的mux的
	})
}

func run() error { //启动 http 服务，绑定 http 和 grpc
	ctx := context.Background() //创建上下文对象
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()                                                             // 创建路由器处理 http 请求
	handlerWithMiddleware := loggingMiddleware(mux)                                          // 假设调用中间件
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}      // 创建时的option 这个也可以变成全局变量
	err := pb.RegisterCreateNftClassHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts) // 创建连接后进行绑定
	if err != nil {
		log.Fatalln(err)
	}
	return http.ListenAndServe(":8081", handlerWithMiddleware) // 启动 http 服务监听端口
}
