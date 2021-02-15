package services

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/xcltapestry/gowk/core/naming"
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


func (s *RPCService) Initialize() (*grpc.Server,error) {
	var err error
	s.listen, err = net.Listen("tcp", s.port)
	if err != nil {
		return nil,fmt.Errorf("failed start server::%s", err)
	}
	s.grpcServer = grpc.NewServer()

	return s.grpcServer,nil
}

func (s *RPCService) Run() error {
	fmt.Println("Register Run()...")
	go func() {
		fmt.Println("Register grpcServer.Serve(listen)...")
		if err := s.grpcServer.Serve(s.listen); err != nil {
			log.Fatal("failed serve:",err)
		}
		fmt.Printf("gRPC服务停止监听端口: %s\n", s.port)
	}()

	return s.serviceRegister()
}


func (s *RPCService) serviceRegister() error{
	s.rpcSvcAddr = fmt.Sprintf("127.0.0.1%s", s.port)
	// fmt.Println("开始注册服务  服务名:", _serviceName, " 端口:", rpcSvcAddr)
	var err error
	s.discovery = naming.NewNaming(naming.WithAddress([]string{"localhost:2379"}))
	err = s.discovery.Register(_serviceName, s.rpcSvcAddr)
	if err != nil {
		return fmt.Errorf("Register err: %s", err.Error())
	}
	fmt.Println("gRPC服务注册成功!")
	return nil 
}

func (s *RPCService) Stop(ctx context.Context) {
	if s.discovery != nil {
		err := s.discovery.Deregister(_serviceName, s.rpcSvcAddr)
		if err != nil {
			fmt.Println(" Deregister err:", err," _serviceName",_serviceName)
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

