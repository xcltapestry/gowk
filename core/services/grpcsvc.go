package services

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
	"net"
	"strings"

	"google.golang.org/grpc"
	"github.com/xcltapestry/gowk/core/naming"
	"github.com/xcltapestry/gowk/pkg/logger"
)

var _serviceName = "_serviceName"

type RPCService struct {
	rpcSvcAddr string

	port       string
	listen     net.Listener
	grpcServer *grpc.Server

	discovery *naming.Naming
}

func NewRPCService() *RPCService {
	svc := &RPCService{}
	svc.port = ":8082"
	return svc
}


// func (s *RPCService) Initialize() (*grpc.Server,error) {
// 	var err error
// 	s.listen, err = net.Listen("tcp", s.port)
// 	if err != nil {
// 		return nil,fmt.Errorf("failed start server::%s", err)
// 	}
// 	s.grpcServer = grpc.NewServer()

// 	return s.grpcServer,nil
// }

func (s *RPCService) Listen(addrs... string) (*grpc.Server,error) {

	var addr string
	for _, a := range addrs {
		addr = a
	}

	if strings.TrimSpace(addr) != "" {
		s.port = addr
	}

	var err error
	s.listen, err = net.Listen("tcp", s.port)
	if err != nil {
		return nil,fmt.Errorf("failed start server:%s port:%s", err,s.port)
	}
	s.grpcServer = grpc.NewServer()

	return s.grpcServer,nil
}

func (s *RPCService) Run() error {
	logger.Info("gRPC服务开始启动...")
	go func() {
		logger.Info("gRPC服务准备监控端口...")
		if err := s.grpcServer.Serve(s.listen); err != nil {
			logger.Fatal("failed serve:",err)
		}
		fmt.Printf("gRPC服务停止监听端口: %s\n", s.port)
	}()

	return s.registerService()
}


func (s *RPCService) registerService() error{
	s.rpcSvcAddr = fmt.Sprintf("127.0.0.1%s", s.port)
	// fmt.Println("开始注册服务  服务名:", _serviceName, " 端口:", rpcSvcAddr)
	var err error
	s.discovery = naming.NewNaming(naming.WithAddress([]string{"localhost:2379"}))
	err = s.discovery.Register(_serviceName, s.rpcSvcAddr)
	if err != nil {
		return fmt.Errorf("Register err: %s", err.Error())
	}
	logger.Info("gRPC服务注册成功!")
	return nil 
}

func (s *RPCService) Stop(ctx context.Context) {
	if s.discovery != nil {
		err := s.discovery.Deregister(_serviceName, s.rpcSvcAddr)
		if err != nil {
			logger.Info(" Deregister err:", err," _serviceName",_serviceName)
		}
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	
	if s.listen != nil {
		s.listen.Close()
	}
}

func (s *RPCService) GetgRPCServer() *grpc.Server {
	return s.grpcServer
}

