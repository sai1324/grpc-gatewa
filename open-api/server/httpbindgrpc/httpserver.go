package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-gateway/open-api/internal/api"
	configInfo "grpc-gateway/open-api/pkg/configs/config"
	"log"
	"net/http"
	"strings"
)

var connections map[string]*grpc.ClientConn

func main() {
	sp := []runtime.ServeMuxOption{
		// 响应前调用
		//runtime.WithForwardResponseOption(middleware.Forward),
	}
	// 创建处理所需的 路由器，处理http请求
	mux := runtime.NewServeMux(sp...)

	// 读取配置项 初始化 gRPC 连接
	initgrpcCon()

	// 自定义路由的方式 ,处理对外的http请求
	err := mux.HandlePath(http.MethodPost, "/nft/classes", handleCreateNftClass)
	if err != nil {
		fmt.Println(err)
		return
	}
	//  自定义路由
	err = mux.HandlePath(http.MethodGet, "/nft/classes/{id}", handleClassByID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 启动 HTTP 服务器 向外暴露http服务 为本机8080端口
	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		return
	}
}

// 初始化 gRPC 连接
func initgrpcCon() {
	viper.SetConfigFile("configs/config.toml")
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var config configInfo.GRPCClientConfig
	err = viper.Unmarshal(&config) // 将配置文件解析到结构体中
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	connections = make(map[string]*grpc.ClientConn)
	for serviceName, clientInfo := range config.Clients { //循环创建连接
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%s", clientInfo.Host, clientInfo.Port),
			grpc.WithInsecure())
		if err != nil {
			log.Fatalf("连接到 %s grpc 失败：%v", serviceName, err)
			return
		}
		connections[serviceName] = conn
		fmt.Sprintf("%s:%s", clientInfo.Host, clientInfo.Port)
	}

}

func handleCreateNftClass(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	var response *api.CreateNftResponse
	//var err error
	// 预先创建好连接
	//var clientShenzhouMap = map[string]api.ShenzhouClient{
	//	"shenzhou_native":      api.NewShenzhouClient(connections["shenzhou_native"]),
	//}
	//var clientWenchangchainMap = map[string]api.WenchangchainClient{
	//	"wenchangchain_native": api.NewWenchangchainClient(connections["wenchangchain_native"]),
	//}

	chain := r.Header.Get("api-key")
	log.Println("api-key 值", chain) // 打印看下 用于调试

	ctx := context.Background() //创建上下文对象
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	//  修改这个if else 这样选择客户端不好，事先创建好连接放在map中比较好
	if strings.ToLower(chain) == "wenchangchain_native" {
		err := api.RegisterWenchangchainHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
		if err != nil {
			return
		}
		mux.ServeHTTP(w, r)

		// 不使用grpc-gateway的方式
		//client := api.NewWenchangchainClient(connections[chain]) //创建对应客户端
		//response, err = client.CreateNftClass(context.Background(), &api.CreateNftRequest{})
		//if err != nil {
		//	// 处理错误
		//	fmt.Println(response, "wenchang---error") //简单输出一下测试
		//}

	} else if strings.ToLower(chain) == "shenzhou_evm" {
		err := api.RegisterShenzhouHandlerFromEndpoint(ctx, mux, "localhost:50052", opts)
		if err != nil {
			return
		}
		mux.ServeHTTP(w, r)

		//client := api.NewShenzhouClient(connections[chain])
		//response, err = client.CreateNftClass(context.Background(), &api.CreateNftRequest{})
		if err != nil {
			// 处理错误
			fmt.Println(response, "shenzhou---error") //简单输出一下测试
		}

	} else {
		fmt.Println("post方法未匹配处理函数")
	}
	//fmt.Println(response)
	//jsonData, err := json.Marshal(response)
	//if err != nil {
	//	// 处理错误
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//// 设置响应头
	//w.Header().Set("Content-Type", "application/json")
	//
	//// 发送 JSON 响应数据给客户端
	//w.Write(jsonData)
}

func handleClassByID(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	chain := r.Header.Get("api-key")
	// 根据字段的值，调用相应的 gRPC 服务
	var response *api.ClassByIDResponse
	var err error
	// 修改这个if else --减小冗余，但是如何创建对应的连接
	if strings.ToLower(chain) == "wenchangchain_native" {
		client := api.NewWenchangchainClient(connections[chain]) //根据传入的连接名创建对应客户端
		response, err = client.ClassByID(context.Background(), &api.ClassByIDRequest{})
		if err != nil {
			fmt.Println(response) //简单输出
			// 处理错误
		}
	} else if strings.ToLower(chain) == "shenzhou_evm" {
		client := api.NewShenzhouClient(connections[chain])
		response, err = client.ClassByID(context.Background(), &api.ClassByIDRequest{})
		if err != nil {
			fmt.Println(response, "shenzhou_native") //简单输出一下测试
		}
		// 处理响应
	} else {
		fmt.Println("get方法未匹配处理函数")
	}
	// 返回 返回体
	jsonData, err := json.Marshal(response)
	if err != nil {
		// 处理错误
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 发送 JSON 响应数据给客户端
	w.Write(jsonData)
}
