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
	"fmt"
	"sync"
	"time"

	etcdClientV3 "github.com/coreos/etcd/clientv3" // etcd的包管理让人心碎
	"google.golang.org/grpc"
)

var _DefaultEndpoints = "localhost:2379"
var _DefaultEtcdPrefix = "/discovery/etcd/grpc/service/"

//Naming 基于ETCD的服务注册与发现
type Naming struct {
	config         etcdClientV3.Config
	prefix         string
	requestTimeout time.Duration

	kvRegister sync.Map
	kvResolver sync.Map
}

func (n *Naming) GetConfig() etcdClientV3.Config {
	return n.config
}

//////////////////////////////////////////////////////////////////
//
// Register
//
//////////////////////////////////////////////////////////////////

//Register 服务注册
func (n *Naming) Register(serviceName, addr string) error {
	var err error
	s := &etcdRegisty{}
	s.ServiceInfo = &ServiceInfo{Name: serviceName, Addr: addr,
		GrpcProxyEndpoint: n.GetGrpcProxyEndpoint(serviceName)}
	err = s.Register(n.config, s.ServiceInfo)
	if err != nil {
		return err
	}

	n.kvRegister.Store(s.GetKey(), s)
	return nil
}

//Deregister 释放
func (n *Naming) Deregister(serviceName, addr string) error {
	s := &etcdRegisty{}
	s.ServiceInfo = &ServiceInfo{Name: serviceName, Addr: addr,
		GrpcProxyEndpoint: n.GetGrpcProxyEndpoint(serviceName)}
	s.Deregister()
	n.kvRegister.Delete(s.GetKey())
	return nil
}

//GetGrpcProxyEndpoint 获取proxy endpoint
func (n *Naming) GetGrpcProxyEndpoint(serviceName string) string {
	// "/discovery/etcd/grpc/service/%s"
	return fmt.Sprintf("%s%s", n.prefix, serviceName)
}

//////////////////////////////////////////////////////////////////
//
// Resolver
//
//////////////////////////////////////////////////////////////////
//GetResolver 命名解析
func (n *Naming) GetResolver(serviceName string) error {
	var err error
	s := &etcdRegisty{}
	s.ServiceInfo = &ServiceInfo{Name: serviceName, Addr: "",
		GrpcProxyEndpoint: n.GetGrpcProxyEndpoint(serviceName)}

	r := &etcdResolver{}
	err = r.GetResolver(n.config, s.ServiceInfo)
	if err != nil {
		return fmt.Errorf(" Resolver失败! err:%s", err)
	}
	n.kvResolver.Store(s.GetKey(), r)
	return nil
}

//GetClientConn 得到grpc客户端连接
func (n *Naming) GetClientConn(serviceName string) (*grpc.ClientConn, error) {

	s := &etcdRegisty{}
	s.ServiceInfo = &ServiceInfo{Name: serviceName, Addr: ""}

	r, ok := n.kvResolver.Load(s.GetKey())
	if !ok {
		return nil, fmt.Errorf(" load(%s) 发生异常!", s.GetKey())
	}

	o, ok := r.(*etcdResolver)
	if !ok {
		return nil, fmt.Errorf(" 转换为Resolver时发生异常!")
	}

	return o.GetClientConn(), nil
}
