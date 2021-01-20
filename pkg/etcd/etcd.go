package etcd

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
 * @Project gowk
 * @Description go framework
 * @author XiongChuanLiang<br/>(xcl_168@aliyun.com)
 * @license http://www.apache.org/licenses/  Apache v2 License
 * @version 1.0
 */
 
import (
	"time"

	"go.etcd.io/etcd/clientv3"
)

type EtcdCli struct {
	config         clientv3.Config
	prefix         string
	requestTimeout time.Duration
	client         *clientv3.Client
}

type EtcdCliOption func(*EtcdCli)

func NewEtcdCli(options ...func(*EtcdCli)) *EtcdCli {

	cli := &EtcdCli{}
	cli.config.Endpoints = []string{"localhost:2379"}
	cli.config.DialTimeout = 2 * time.Second

	for _, f := range options {
		f(cli)
	}

	return cli
}

func WithAddress(addrs []string) EtcdCliOption {
	return func(c *EtcdCli) {
		c.config.Endpoints = addrs
	}
}

func WithDialTimeout(dialTimeout time.Duration) EtcdCliOption {
	return func(c *EtcdCli) {
		c.config.DialTimeout = dialTimeout
	}
}

func WithPrefix(prefix string) EtcdCliOption {
	return func(c *EtcdCli) {
		c.prefix = prefix
	}
}

func WithRequestTimeout(requestTimeout time.Duration) EtcdCliOption {
	return func(c *EtcdCli) {
		c.requestTimeout = requestTimeout
	}
}

func (e *EtcdCli) Connect() error {
	var err error
	e.client, err = clientv3.New(e.config)
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdCli) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *EtcdCli) Client() *clientv3.Client {
	return e.client
}

func (e *EtcdCli) GetPrefix() string {
	return e.prefix
}

func (e *EtcdCli) GetRequestTimeout() time.Duration {
	return e.requestTimeout
}
