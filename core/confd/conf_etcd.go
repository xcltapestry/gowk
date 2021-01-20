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
 * @Project gowk
 * @Description go framework
 * @author XiongChuanLiang<br/>(xcl_168@aliyun.com)
 * @license http://www.apache.org/licenses/  Apache v2 License
 * @version 1.0
 */

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
	"github.com/xcltapestry/gowk/pkg/etcd"
)

type CfdEtcd struct {
	currentAppConfig *AppConfig
	watchOnce        sync.Once
}

func NewCfdEtcd(appConfig *AppConfig) *CfdEtcd {
	ce := &CfdEtcd{}
	ce.currentAppConfig = appConfig
	return ce
}

func (ce *CfdEtcd) newEtcdCli(appConfig *AppConfig) *etcd.EtcdCli {
	return etcd.NewEtcdCli(
		etcd.WithAddress(appConfig.GetConfdRemoteAddrs()), //"localhost:2379"
		etcd.WithDialTimeout(2*time.Second),
		etcd.WithRequestTimeout(1*time.Second),
		etcd.WithPrefix(appConfig.GetRootKey()))
}

//loadConfigFromRemote 依server的rootkey，取出对应的配置文件，并解析更新当前配置
func (ce *CfdEtcd) LoadConfigFromRemote(loadViper *viper.Viper, rootKey, configType string) error {

	if loadViper == nil {
		return fmt.Errorf(" viper is null.")
	}

	cli := ce.newEtcdCli(ce.currentAppConfig)
	err := cli.Connect()
	if err != nil {
		return err
	}
	defer cli.Close()

	gresp, err := cli.Client().Get(context.Background(), rootKey, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	loadViper.SetConfigType(configType)

	var configData *bytes.Buffer
	configData = bytes.NewBuffer(gresp.Kvs[0].Value)
	err = loadViper.ReadConfig(configData)
	if err != nil {
		return err
	}

	return err
}

//WatchRemoteConfig 监控ETCD,保存配置热更新,不建议在产线使用。产线配置变更建议通过重新发布蓝绿部署等方式，避免产线事故发生
func (ce *CfdEtcd) WatchRemoteConfig(loadViper *viper.Viper, rootKey, configType string) error {

	cli := ce.newEtcdCli(ce.currentAppConfig)
	err := cli.Connect()
	if err != nil {
		return err
	}
	defer cli.Close()

	ce.watchOnce.Do(func() {
		for {
			rch := cli.Client().Watch(context.Background(), rootKey, clientv3.WithPrefix())
			for wresp := range rch {
				for _, ev := range wresp.Events {
					switch ev.Type {
					case clientv3.EventTypePut:
						fmt.Printf("[Watch]:EventTypePut: [%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

						err2 := ce.LoadConfigFromRemote(loadViper, rootKey, configType)
						if err2 != nil {
							fmt.Println("[ERROR] 同步失败!")
							continue
						}

					case clientv3.EventTypeDelete:
						fmt.Printf("[Watch]:EventTypeDelete: [%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

					}

				}
			}
		} //watch
	})

	return nil
}

func (ce *CfdEtcd) ReadConfigFileToETCD(confFile, rootKey string) error {
	if IsNotExist(confFile) {
		return fmt.Errorf(" 文件不存在! ")
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

	rootKey, value1 := rootKey, string(fileBody)
	if _, err := cli.Client().Put(context.Background(), rootKey, value1); err != nil {
		return fmt.Errorf("err:%s rootKey:%s", err.Error(), rootKey)
	}

	return nil
}
