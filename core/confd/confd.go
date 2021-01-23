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
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/xcltapestry/gowk/pkg/logger"
)

type Confd struct {
	//syncMode 是否需要启用热更新
	syncMode bool // 是否启用同步更新热更新配置的方式？

	appConfig   *AppConfig
	remoteConfd *EtcdConfd
	localConfd  *ConfLocalFile
	viper       *viper.Viper

	rmu sync.RWMutex
}

func NewConfd(appConfig *AppConfig) *Confd {
	c := &Confd{}
	c.appConfig = appConfig
	c.remoteConfd = NewEtcdConfd(appConfig)
	c.localConfd = NewConfLocalFile()
	c.viper = viper.New()
	return c
}

func (cfd *Confd) BindConfig() error {

	if cfd.remoteConfd == nil || cfd.appConfig == nil || cfd.viper == nil {
		return fmt.Errorf(" 对象未初始化。")
	}

	addrs := cfd.appConfig.GetConfdRemoteAddrs()
	confFile := cfd.appConfig.GetLocalConfigFile()

	switch {
	case len(addrs) > 0:
		rKey := cfd.appConfig.GetRootKey()
		logger.Infow("准备从远程ETCD配置中心读取配置. ", " key:", rKey)
		etcdCli := NewEtcdConfd(cfd.appConfig)
		loadViper := viper.New()
		err := etcdCli.LoadConfigFromRemote(loadViper, rKey, cfd.appConfig.GetRemoteConfigType())
		if err != nil {
			return err // todo: 假如err是因为rootkey在远程没找到，则尝试读本地文件
		}
		cfd.rmu.Lock()
		cfd.viper = loadViper
		cfd.rmu.Unlock()
		logger.Infow("从远程ETCD配置中心读取配置完毕. ", " key:", rKey)
		cfd.WatchRemoteConfig()
		return nil
	case confFile != "": //读本地文件
		logger.Infow("准备从本地配置文件读取配置.", "文件", confFile)
		loadViper := viper.New()
		err := cfd.localConfd.LoadConfigFromLocalFile(loadViper, confFile)
		if err != nil {
			return err
		}
		cfd.rmu.Lock()
		cfd.viper = loadViper
		cfd.rmu.Unlock()
		logger.Info("从本地配置文件读取配置完毕. ", "文件", confFile)
		return nil
	default: //len(addrs) == 0 && confFile == ""
		return fmt.Errorf(" 远程配置中心或本地配置文件必须指定一个。")
	}

	return nil
}

func (cfd *Confd) WatchRemoteConfig() error {
	if cfd.appConfig.IsSync() == false {
		return nil
	}
	addrs := cfd.appConfig.GetConfdRemoteAddrs()
	if len(addrs) > 0 {
		etcdCli := NewEtcdConfd(cfd.appConfig)
		loadViper := viper.New()
		err := etcdCli.WatchRemoteConfig(loadViper, cfd.appConfig.GetRootKey(), cfd.appConfig.GetRemoteConfigType())
		if err != nil {
			return err
		}
		cfd.rmu.Lock()
		cfd.viper = loadViper
		cfd.rmu.Unlock()
	}
	return nil
}

func (cfd *Confd) GetRootKey() string {
	return cfd.appConfig.GetRootKey()
}

func (cfd *Confd) Get(key string) interface{} {
	if cfd.viper == nil {
		return nil
	}
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.Get(key)
}

func (cfd *Confd) GetString(key string) string {
	if cfd.viper == nil {
		return ""
	}
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetString(key)
}

func (cfd *Confd) GetBool(key string) bool {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetBool(key)
}

func (cfd *Confd) GetDuration(key string) time.Duration {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetDuration(key)
}

func (cfd *Confd) GetInt(key string) int {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetInt(key)
}

func (cfd *Confd) GetStringMap(key string) map[string]interface{} {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetStringMap(key)
}

func (cfd *Confd) GetStringMapString(key string) map[string]string {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetStringMapString(key)
}

func (cfd *Confd) GetStringMapStringSlice(key string) map[string][]string {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetStringMapStringSlice(key)
}

func (cfd *Confd) GetStringSlice(key string) []string {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.GetStringSlice(key)
}

func (cfd *Confd) InConfig(key string) bool {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.InConfig(key)
}

func (cfd *Confd) IsSet(key string) bool {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.IsSet(key)
}

func (cfd *Confd) AllKeys() []string {
	cfd.rmu.Lock()
	defer cfd.rmu.Unlock()
	return cfd.viper.AllKeys()
}

func (cfd *Confd) ReadConfigFileToRemote() error {
	return cfd.remoteConfd.ReadConfigFileToETCD(cfd.appConfig.GetLocalConfigFile(), cfd.appConfig.GetRootKey())
}

func (cfd *Confd) String() string {
	return cfd.appConfig.String()
}
