package main

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
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/xcltapestry/gowk/_example/protocol"
	"github.com/xcltapestry/gowk/core/app"
	"github.com/xcltapestry/gowk/core/naming"
	"github.com/xcltapestry/gowk/core/services"
	"github.com/xcltapestry/gowk/pkg/logger"
)

var (
	// go build 时通过 ldflags 设置
	codeBuiildSourceVersion = ""
	codeBuildTime           = ""
	gitHash                 = ""
)

/*
test:
	go run service.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -confd.remote.addrs="localhost:2379"

	go run service.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -confd.remote.addrs=""

	go run main.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -logger.output=alsologtostdout

*/

func main() {
	logger.NewDefaultLogger()

	app.New()

	// http 服务
	httpSvc := services.NewHTTPService()
	httpSvc.Router(RegisterHandlers)
	httpSvc.Listen(app.Confd().GetString("application.listen.http"))
	app.Serve(httpSvc)

	// gRPC 服务
	rpcSvc := CreateRPC()
	// 注册 hello RPC服务
	RegisterHelloServer(rpcSvc)
	app.Serve(rpcSvc)

	app.Run()
}

func CreateRPC() *services.RPCService {
	// rpc服务
	rpcSvc := services.NewRPCService()
	_, err := rpcSvc.Listen()
	if err != nil {
		logger.Panic(err)
	}
	// 从配置中心获取配置
	etcdAddrs := app.Confd().GetStringSlice("application.registry.etcdv3.endpoints")
	etcdDialTimeout := app.Confd().GetDuration("application.registry.etcdv3.dialtimeout")
	// 基于etcd的RPC名字服务
	rpcSvc.NewNaming(naming.WithAddress(etcdAddrs), naming.WithDialTimeout(etcdDialTimeout))
	return rpcSvc
}

func RegisterHelloServer(rpcSvc *services.RPCService) {
	// hello 服务
	if err := rpcSvc.Registry("HelloService"); err != nil {
		logger.Panic(err)
	}
	// 注册Hello服务
	pb.RegisterHelloServiceServer(rpcSvc.GetgRPCServer(), &HelloServer{})
}

//RegisterHandlers 路由
func RegisterHandlers(m *mux.Router) {
	m.HandleFunc("/health", HealthHandler)
	m.Handle("/", http.NotFoundHandler())
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug(" Health => ", time.Now().Unix())
	fmt.Fprintf(w, "Health: %v\n", time.Now().Unix())
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

// func displayConfd(){

// 	fmt.Println(" ------------------------------------- ")
// 	keys := app.Confd().AllKeys()
// 	for _, k := range keys {
// 		fmt.Printf("AllKeys----[keys] %s\n", k)
// 	}

// 	logger.Info("server.http: ",app.Confd().GetString("application.listen.http"))
// 	logger.Info("server.grpc: ",app.Confd().GetString("application.listen.grpc"))
// 	logger.Info("etcdv3.endpoints: ",app.Confd().GetStringSlice("application.registry.etcdv3.endpoints"))
// 	logger.Info("etcdv3.connecttimeout: ",app.Confd().GetDuration("application.registry.etcdv3.connecttimeout"))

// }
