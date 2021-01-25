package confd

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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/viper"
	"github.com/xcltapestry/gowk/pkg/etcd"
	"github.com/xcltapestry/gowk/pkg/utils"
	"go.etcd.io/etcd/clientv3"
)

// EtcdConfd
type EtcdConfd struct {
	currentAppConfig *AppConfig
}

func NewEtcdConfd(appConfig *AppConfig) *EtcdConfd {
	ce := &EtcdConfd{}
	ce.currentAppConfig = appConfig
	return ce
}

func (ce *EtcdConfd) newEtcdCli(appConfig *AppConfig) *etcd.EtcdCli {
	return etcd.NewEtcdCli(
		etcd.WithAddress(appConfig.GetConfdRemoteAddrs()), //"localhost:2379"
		etcd.WithDialTimeout(2*time.Second),
		etcd.WithRequestTimeout(1*time.Second),
		etcd.WithPrefix(appConfig.GetRootKey()))
}

//loadConfigFromRemote 依server的rootkey，取出对应的配置文件，并解析更新当前配置
func (ce *EtcdConfd) LoadConfigFromRemote(loadViper *viper.Viper, rootKey, configType string) error {

	if loadViper == nil {
		return fmt.Errorf(" viper is null.")
	}

	cli := ce.newEtcdCli(ce.currentAppConfig)
	err := cli.Connect()
	if err != nil {
		return fmt.Errorf(" 连接ETCD配置中心时，发生异常. err:%s key:%s", err.Error(), rootKey)
	}
	defer cli.Close()

	fmt.Println("[Resolver] connect etcd!")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	gresp, err := cli.Client().Get(ctx, rootKey, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	loadViper.SetConfigType(configType)

	if gresp == nil || len(gresp.Kvs) <= 0 {
		return fmt.Errorf(" 没有找到对应的配置信息. key:%s", rootKey)
	}

	var configData *bytes.Buffer
	configData = bytes.NewBuffer(gresp.Kvs[0].Value)
	err = loadViper.ReadConfig(configData)
	if err != nil {
		return fmt.Errorf(" 读取配置信息时发生异常. err:%s, key:%s len:%d ", err.Error(), rootKey, configData.Len())
	}

	return nil
}

//WatchRemoteConfig 监控ETCD,保存配置热更新,不建议在产线使用。产线配置变更建议通过重新发布蓝绿部署等方式，避免产线事故发生
func (ce *EtcdConfd) WatchRemoteConfig(loadViper *viper.Viper, rootKey, configType string) error {

	if loadViper == nil {
		return fmt.Errorf(" viper is null.")
	}

	cli := ce.newEtcdCli(ce.currentAppConfig)
	err := cli.Connect()
	if err != nil {
		return err
	}
	defer cli.Close()

	rch := cli.Client().Watch(context.Background(), rootKey, clientv3.WithPrefix())
	go func() {
		for {

			for wresp := range rch {
				if err := wresp.Err(); err != nil {
					fmt.Println("获得的Key变更信息有异常. ", err)
					break
				}

				for _, ev := range wresp.Events {
					if ev.Kv == nil {
						continue
					}
					switch ev.Type {
					case clientv3.EventTypePut:
						fmt.Printf("[Watch]:EventTypePut: [%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

						err2 := ce.LoadConfigFromRemote(loadViper, rootKey, configType)
						if err2 != nil {
							fmt.Println("[ERROR] 同步失败! err:", err)
							continue
						}

					case clientv3.EventTypeDelete:
						fmt.Printf("[Watch]:EventTypeDelete: [%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

					} // end switch ev.Type {

				} //end range wresp.Events {
			} //end range rch
			time.Sleep(1 * time.Second)
			rch = cli.Client().Watch(context.Background(), rootKey, clientv3.WithPrefix())

		} // end for{}
	}() // end go func() {

	return nil
}

func (ce *EtcdConfd) ReadConfigFileToETCD(confFile, rootKey string) error {
	if utils.IsNotExist(confFile) {
		return fmt.Errorf(" 文件(%s)不存在! ", confFile)
	}

	cli := ce.newEtcdCli(ce.currentAppConfig)
	err := cli.Connect()
	if err != nil {
		return err
	}
	defer cli.Close()

	fileBody, err := ioutil.ReadFile(confFile)
	if err != nil {
		return fmt.Errorf("err:%s confFile:%s", err.Error(), confFile)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if _, err := cli.Client().Put(ctx, rootKey, string(fileBody)); err != nil {
		return fmt.Errorf("err:%s rootKey:%s", err.Error(), rootKey)
	}

	return nil
}

func (ce *EtcdConfd) SyncToEtcd(addrs, rootkey, confFile string) error {

	// etcdAddrs := strings.Split(addrs, ";")
	// cli := etcd.NewEtcdCli(
	// 			etcd.WithAddress(etcdAddrs), //"localhost:2379"
	// 			etcd.WithDialTimeout(2*time.Second),
	// 			etcd.WithRequestTimeout(1*time.Second),
	// 			etcd.WithPrefix(rootkey))
	// 	}

	//return  cli.ReadConfigFileToETCD(confFile, rootKey)

	return nil
}
