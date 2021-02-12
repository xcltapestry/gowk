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

/**

            +-------------+
            | etcd server |
            +------+------+
                   ^ watch key A (s-watcher)
                   |
           +-------+-----+
           | gRPC proxy  | <-------+
           |             |         |
           ++-----+------+         |watch key A (c-watcher)
watch key A ^     ^ watch key A    |
(c-watcher) |     | (c-watcher)    |
    +-------+-+  ++--------+  +----+----+
    |  client |  |  client |  |  client |
    |         |  |         |  |         |
    +---------+  +---------+  +---------+

https://etcd.io/docs/v3.1.12/op-guide/grpc_proxy/

**/

import (
	"fmt"

	etcdClientV3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/proxy/grpcproxy"
)

const (
	_defaultTTL int = 60
)

type ServiceInfo struct {
	Name              string //服务名
	Addr              string //服务地址
	GrpcProxyEndpoint string //etcd注册endpoint
}

type etcdRegisty struct {
	ServiceInfo *ServiceInfo
	EtcdClient  *etcdClientV3.Client
}

//GetKey 得到服务Key
func (svc *etcdRegisty) GetKey() string {
	return fmt.Sprint(svc.ServiceInfo.Name, "/", svc.ServiceInfo.Addr)
}

//Register 注册服务
func (s *etcdRegisty) Register(etcdCfg etcdClientV3.Config, rpcSvc *ServiceInfo, ttl int) error {
	if rpcSvc == nil {
		return fmt.Errorf(" Register失败. rpcSvc is nill.")
	}
	if ttl <= 0 {
		ttl = _defaultTTL
	}

	var err error
	s.EtcdClient, err = etcdClientV3.New(etcdCfg)
	if err != nil {
		return fmt.Errorf(" ETCD客户端连接失败. err:%s", err)
	}

	grpcproxy.Register(s.EtcdClient, rpcSvc.GrpcProxyEndpoint, rpcSvc.Addr, ttl)
	return nil
}

//Deregister 释放
func (s *etcdRegisty) Deregister() error {
	if s.EtcdClient != nil {
		s.EtcdClient.Close()
	}
	return nil
}
