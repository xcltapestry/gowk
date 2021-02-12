package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xcltapestry/gowk/core/naming"

	pb "github.com/xcltapestry/gowk/_example/protocol"
)

var _serviceName string = "hellogrpcsvc2"
var _etcdAddr string = "localhost:2379"

func main() {
	discovery := naming.NewNaming()
	err := discovery.GetResolver(_serviceName)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		request := &pb.HelloRequest{Greeting: fmt.Sprintf("%d", i)}
		//从etcd获取rpc连接
		rpcConn, err := discovery.GetClientConn(_serviceName)
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
