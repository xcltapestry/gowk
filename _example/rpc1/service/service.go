package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/xcltapestry/gowk/core/naming"
	"google.golang.org/grpc"

	pb "github.com/xcltapestry/gowk/_example/protocol"
)

var _serviceName string = "hellogrpcsvc2"
var _etcdAddr string = "localhost:2379"

func main() {
	port := ":8082"
	//rpc服务监听指定端口
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("rpc服务监听指定端口: ", port)

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &HelloServer{})
	go func() {
		// time.Sleep(500 * time.Millisecond)
		fmt.Println("Register grpcServer.Serve(listen)...")
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("gRPC服务停止监听端口: %s\n", port)
	}()

	rpcSvcAddr := fmt.Sprintf("127.0.0.1%s", port)
	fmt.Println("开始注册服务  服务名:", _serviceName, " 端口:", rpcSvcAddr)
	// --  Register

	discovery := naming.NewNaming(naming.WithAddress([]string{"localhost:2379"}))
	err = discovery.Register(_serviceName, rpcSvcAddr)
	if err != nil {
		log.Fatal("Register err: ", err.Error())
	}
	fmt.Println("gRPC服务注册成功!")

	// 监听信息量
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sign := <-sigs
	fmt.Printf("接收到信号: %v, gRPC服务将关闭。\n", sign)

	// -- 清理
	err = discovery.Deregister(_serviceName, rpcSvcAddr)
	if err != nil {
		fmt.Println(" UnRegister v1  err:", err)
	}

	grpcServer.GracefulStop()
	listen.Close()
}

//HelloServer type
type HelloServer struct{}

//SayHello func
func (h *HelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	response := &pb.HelloResponse{
		Reply: fmt.Sprintf("hello, %s", req.Greeting),
	}
	fmt.Println(req.Greeting)
	return response, nil
}
