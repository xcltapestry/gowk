package naming

/**
 * Copyright 2021  gowrk Author. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"context"
	"fmt"
	"math"
	"time"

	etcdClientV3 "github.com/coreos/etcd/clientv3"
	etcdNaming "github.com/coreos/etcd/clientv3/naming"
	"google.golang.org/grpc"
)

const (
	//google.golang.org/grpc/server.go.defaultServerOptions
	_WriteBufferSize       int   = 256 * 1024 //32kb
	_ReadBufferSize        int   = 256 * 1024 //32kb
	_InitialWindowSize     int32 = 1 << 30
	_InitialConnWindowSize int32 = 1 << 30
)

type etcdResolver struct {
	grpcConn    *grpc.ClientConn
	serviceName string
}

//GetResolver gRPC客户端
func (r *etcdResolver) GetResolver(etcdCfg etcdClientV3.Config, rpcSvc *ServiceInfo) error {
	if rpcSvc == nil {
		return fmt.Errorf(" Resolver失败. rpcSvc is nill.")
	}

	cli, err := etcdClientV3.New(etcdCfg)
	if err != nil {
		return fmt.Errorf(" ETCD连接失败. err:%s", err)
	}
	defer cli.Close()

	grpcResolver := &etcdNaming.GRPCResolver{Client: cli}
	lb := grpc.RoundRobin(grpcResolver) //轮询方式
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	gconn, gerr := grpc.DialContext(ctx, rpcSvc.GrpcProxyEndpoint,
		grpc.WithInsecure(),
		grpc.WithBalancer(lb),
		grpc.WithTimeout(time.Second*10),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.FailFast(false),
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithBackoffConfig(grpc.BackoffConfig{
			MaxDelay: time.Second * 10,
		}),
		grpc.WithWriteBufferSize(_WriteBufferSize),
		grpc.WithReadBufferSize(_ReadBufferSize),
		grpc.WithInitialWindowSize(_InitialWindowSize),
		grpc.WithInitialConnWindowSize(_InitialConnWindowSize),
	)
	if gerr != nil {
		return fmt.Errorf(" 连接etcd resolver server 发生异常。 服务名:%s err:%s",
			rpcSvc.GrpcProxyEndpoint, gerr.Error())
	}
	r.grpcConn = gconn
	return nil
}

//GetClientConn 获取gRPC连接
func (r *etcdResolver) GetClientConn() *grpc.ClientConn {
	return r.grpcConn
}
