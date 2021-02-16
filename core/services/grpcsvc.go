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

type RPCService struct {
	gRPCPort       string
	listen     net.Listener
	grpcServer *grpc.Server
	discovery *naming.Naming
	rpcServices map[string]interface{} 
}

func NewRPCService() *RPCService {
	svc := &RPCService{}
	svc.gRPCPort = ":9099"
	svc.rpcServices = make(map[string]interface{})
	return svc
}

func (s *RPCService) Listen(addrs... string) (*grpc.Server,error) {

	var addr string
	for _, a := range addrs {
		addr = a
	}

	if strings.TrimSpace(addr) != "" {
		s.gRPCPort = addr
	}

	var err error
	s.listen, err = net.Listen("tcp", s.gRPCPort)
	if err != nil {
		return nil,fmt.Errorf("failed start server:%s port:%s", err,s.gRPCPort)
	}
	s.grpcServer = grpc.NewServer()

	return s.grpcServer,nil
}

func (s *RPCService) NewNaming(options ...func(*naming.Naming)) error {
	s.discovery = naming.NewNaming(options...)
	return nil 
}

func (s *RPCService) Registry(rpcName string)error{
	if s.discovery  == nil {
		return fmt.Errorf(" 需先执行NewNaming()，再执行Registry().")
	}
	var err error
	err = s.discovery.Register(rpcName, s.gRPCPort)
	if err != nil {
		return fmt.Errorf("Register err: %s", err.Error())
	}
	s.rpcServices[rpcName] = "" 
	logger.Info("RPC服务(",rpcName,")注册成功!")
	return nil 
}


func (s *RPCService) Run() error {
	logger.Info("RPC服务开始启动...")
	go func() {
		logger.Info("RPC服务准备监控端口...")
		if err := s.grpcServer.Serve(s.listen); err != nil {
			logger.Panic("failed serve:",err)
		}
		logger.Infof("RPC服务停止监听端口: %s\n", s.gRPCPort)
	}()
	return nil 
}

func (s *RPCService) Stop(ctx context.Context) {
	if s.discovery != nil {
		for k,_ := range  s.rpcServices {
			err := s.discovery.Deregister(k, s.gRPCPort)
			if err != nil {
				logger.Error(" Deregister err:", err," service:",k)
			}else{
				logger.Info(" Deregister service:",k)
			}
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






