package naming


import (
	"fmt"
	"context"
	"time"

	etcdClientV3 "github.com/coreos/etcd/clientv3"
	etcdNaming "github.com/coreos/etcd/clientv3/naming"
	"google.golang.org/grpc"
)


type etcdResolver struct{
	grpcConn    *grpc.ClientConn
	serviceName string
}

func (e *etcdResolver) connect(addrs string) error {
	config := etcdClientV3.Config{
		Endpoints:   []string{"localhost:2379"},  //"localhost:2379"
		DialTimeout: 3 * time.Second,
	}
	cli, err := etcdClientV3.New(config)
	if err != nil {
		return fmt.Errorf(" err:%s",err)
	}
	defer cli.Close()

	grpcResolver := &etcdNaming.GRPCResolver{Client: cli}
	// pickfirstBalancer: 只使用一个服务地址 
	// RoundRobin : 轮询调度
	// grpclb: 使用一个单独的服务提供负载均衡信息
	lb := grpc.RoundRobin(grpcResolver) 

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	conn, gerr := grpc.DialContext(ctx,e.serviceName,
		grpc.WithInsecure(),
		grpc.WithBalancer(lb),
		grpc.WithTimeout(time.Second*5),
		grpc.WithBlock(),
	)
	if gerr != nil {
		return fmt.Errorf(" 连接etcd resolver server 发生异常。 服务名:%s err:%s",e.serviceName, gerr.Error())
	}
	e.grpcConn = conn
	return nil 
}

func (e *etcdResolver) getGrpcConnection() *grpc.ClientConn {
	return e.grpcConn
}



