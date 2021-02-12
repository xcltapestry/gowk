package naming

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "github.com/xcltapestry/gowk/_example/protocol"
)

var _serviceName string = "hellogrpcsvc2"
var _etcdAddr string = "localhost:2379"

func TestRun(t *testing.T) {
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("启动RPC客户端")
		rpcClient()
	}()
	fmt.Println("启动RPC服务端")
	rpcSvc()
	time.Sleep(9 * time.Second)
}

func rpcSvc() {
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
		fmt.Printf("gRPC服务已启动，监听端口: %s\n", port)
	}()

	rpcSvcAddr := fmt.Sprintf("127.0.0.1%s", port)
	fmt.Println("开始注册服务  服务名:", _serviceName, " 端口:", rpcSvcAddr)
	// --  Register

	regSvc := NewNaming()
	err = regSvc.Register(_serviceName, rpcSvcAddr)
	if err != nil {
		log.Fatal("Register err: ", err.Error())
	}
	fmt.Println("gRPC服务注册成功!")

	err = regSvc.Register(_serviceName+"v2", rpcSvcAddr)
	if err != nil {
		log.Fatal("Register v2 err: ", err.Error())
	}
	fmt.Println("gRPC服务注册成功! v2")

	// 监听信息量
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sign := <-sigs
	fmt.Printf("接收到信号: %v, gRPC服务将关闭。\n", sign)

	// -- 清理
	err = regSvc.Deregister(_serviceName, rpcSvcAddr)
	if err != nil {
		fmt.Println(" UnRegister v1  err:", err)
	}

	err = regSvc.Deregister(_serviceName+"v2", rpcSvcAddr)
	if err != nil {
		fmt.Println(" UnRegister v2  err:", err)
	}

	grpcServer.GracefulStop()
	listen.Close()
}

func rpcClient() {
	namingSvc := NewNaming()
	err := namingSvc.GetResolver(_serviceName)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		request := &pb.HelloRequest{Greeting: fmt.Sprintf("%d", i)}
		//从etcd获取rpc连接
		rpcConn, err := namingSvc.GetClientConn(_serviceName)
		if err != nil {
			log.Fatal(err)
		}
		client := pb.NewHelloServiceClient(rpcConn)
		//发送信息
		resp, err := client.SayHello(context.Background(), request)
		if err != nil {
			fmt.Println(" err:", err)
		} else {
			fmt.Println(" resp:", resp)
		}
		fmt.Printf("test => resp: %+v, err: %+v\n", resp, err)
		time.Sleep(time.Second)
	}

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
