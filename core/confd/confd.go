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
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Confd struct {
	//syncMode 是否需要启用热更新
	syncMode bool // 是否启用同步更新热更新配置的方式？

	appConfig   *AppConfig
	remoteConfd *CfdEtcd
	localConfd  *ConfLocalFile
	viper       *viper.Viper

	rmu       sync.RWMutex
	watchOnce sync.Once
}

func NewConfd(appConfig *AppConfig) *Confd {
	c := &Confd{}
	c.appConfig = appConfig
	c.remoteConfd = NewCfdEtcd(appConfig)
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
		etcdCli := NewCfdEtcd(cfd.appConfig)
		loadViper := viper.New()
		err := etcdCli.LoadConfigFromRemote(loadViper, cfd.appConfig.GetRootKey(), cfd.appConfig.GetRemoteConfigType())
		if err != nil {
			return err
		}
		cfd.rmu.Lock()
		defer cfd.rmu.Unlock()
		cfd.viper = loadViper

		cfd.WatchRemoteConfig()
		return nil
	case confFile != "": //读本地文件
		loadViper := viper.New()
		err := cfd.localConfd.LoadConfigFromLocalFile(loadViper, confFile)
		if err != nil {
			return err
		}
		cfd.rmu.Lock()
		defer cfd.rmu.Unlock()
		cfd.viper = loadViper
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
		etcdCli := NewCfdEtcd(cfd.appConfig)
		loadViper := viper.New()
		err := etcdCli.WatchRemoteConfig(loadViper, cfd.appConfig.GetRootKey(), cfd.appConfig.GetRemoteConfigType())
		if err != nil {
			return err
		}
		cfd.rmu.Lock()
		defer cfd.rmu.Unlock()
		cfd.viper = loadViper
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

// func (cfd *Confd) ListKeys() []string {

// 	fmt.Println(" ------------------------------------- ")
// 	keys := viper.AllKeys()
// 	for _, k := range keys {
// 		fmt.Printf("AllKeys----[keys] %s\n", k)
// 	}
// 	return keys
// }

// func (cfd *Confd) ReadInConfig(path, name string) error {

// 	if cfd.viper == nil {
// 		cfd.viper = viper.New()
// 		cfd.viper.SetConfigType(_defaultConfigType)
// 	}
// 	cfd.viper.SetConfigName(name)
// 	cfd.viper.AddConfigPath(path)

// 	// if err := cfd.viper.ReadInConfig(); err != nil {
// 	// 	// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 	// 	// }
// 	// 	return err
// 	// }

// 	cfd.rwm.Lock()
// 	defer cfd.rwm.Unlock()

// 	return cfd.viper.ReadInConfig()
// }

// func (cfd *Confd) loadConfigFromLocalFile() error {

// 	confFile := cfd.appConfig.GetLocalConfigFile()
// 	path, ext, err := PaseConfigFile(confFile)
// 	if err != nil {
// 		return err
// 	}

// 	loadViper := viper.New()
// 	loadViper.AddConfigPath(path)
// 	path2, _ := os.Getwd()
// 	loadViper.AddConfigPath(path2)

// 	switch ext {
// 	case "json", "hcl", "prop", "props", "properties", "dotenv", "env", "toml", "yaml", "yml", "ini":
// 		loadViper.SetConfigType(ext)
// 	default:
// 		loadViper.SetConfigType("yaml")
// 	}

// 	//读取配置
// 	err = loadViper.ReadInConfig()
// 	if err != nil {
// 		return err
// 	}

// 	cfd.rwm.Lock()
// 	defer cfd.rwm.Unlock()
// 	//更新为新的配置项
// 	cfd.viper = loadViper

// 	return nil
// }

// func (cfd *Config) ReadConfig(in io.Reader, configType string) error {
// 	loadViper := viper.New()
// 	switch configType {
// 	case "json", "hcl", "prop", "props", "properties", "dotenv", "env", "toml", "yaml", "yml", "ini":
// 		loadViper.SetConfigType(configType)
// 	default:
// 		loadViper.SetConfigType("yaml")
// 	}

// 	err := loadViper.ReadConfig(in)
// 	if err != nil {
// 		return err
// 	}

// 	cfd.rwm.Lock()
// 	defer cfd.rwm.Unlock()
// 	//更新为新的配置项
// 	cfd.viper = loadViper

// 	return nil
// }

// //loadConfigFromRemote 依server的rootkey，取出对应的配置文件，并解析更新当前配置
// func (cfd *Confd) LoadConfigFromRemote() error {

// 	cfg.mu.RLock()
// 	defer cfg.mu.RUnlock()

// 	gresp, err := cfg.etcdCli.Client().Get(context.Background(), cfd.appConfig.GetRootKey(), clientv3.WithPrefix())
// 	if err != nil {
// 		return err
// 	}

// 	confFile := cfd.appConfig.GetLocalConfigFile()

// 	path, ext, err := PaseConfigFile(confFile)
// 	if err != nil {
// 		return err
// 	}

// 	viper.AddConfigPath(path)
// 	path2, _ := os.Getwd()
// 	viper.AddConfigPath(path2)

// 	switch ext {
// 	case "yaml", "ini", "toml":
// 		viper.SetConfigType(ext)
// 	default:
// 		viper.SetConfigType("yaml")
// 	}

// 	var configData *bytes.Buffer
// 	configData = bytes.NewBuffer(gresp.Kvs[0].Value)
// 	err = viper.ReadConfig(configData)

// 	xx := viper.Get("Service.Name")
// 	fmt.Println(" Service.Name", xx)
// 	fmt.Println(" namespace:", viper.GetString("Namespace"))

// 	return err
// }

// func (cfd *Confd) ConnectEtcd() error {

// 	cfg.etcdCli = etcd.NewEtcdCli(
// 		etcd.WithAddress(cfd.appConfig.GetConfdRemoteAddrs()), //"localhost:2379"
// 		etcd.WithDialTimeout(2*time.Second),
// 		etcd.WithRequestTimeout(1*time.Second),
// 		etcd.WithPrefix(cfd.appConfig.GetRootKey()))
// 	err := cfg.etcdCli.Connect()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (cfd *Confd) Close() {
// 	if cfg.etcdCli != nil {
// 		cfg.etcdCli.Close()
// 	}
// }

// func (cfd *Confd) ReadConfigFileToETCD2() error {

// 	confFile := cfd.appConfig.GetLocalConfigFile()

// 	if IsNotExist(confFile) {
// 		return fmt.Errorf(" 文件不存在! ")
// 	}

// 	fileBody, err := ioutil.ReadFile(confFile)
// 	if err != nil {
// 		return fmt.Errorf("err:%s confFile:%s", err.Error(), confFile)
// 	}

// 	rootKey, value1 := cfd.appConfig.GetRootKey(), string(fileBody)
// 	ctx, cancel := context.WithTimeout(context.Background(), cfg.etcdCli.GetRequestTimeout())
// 	defer cancel()
// 	if _, err := cfg.etcdCli.Client().Put(ctx, rootKey, value1); err != nil {
// 		return fmt.Errorf("err:%s rootKey:%s", err.Error(), rootKey)
// 	}

// 	return nil

// }

// //ReadConfigFileToETCD 读取yaml之类配置文件，保存至etcd
// func (cfd *Confd) ReadConfigFileToETCD(fileName string) error {
// 	fmt.Println("[ReadfileToETCD] ---")

// 	if strings.TrimSpace(fileName) == "" {
// 		return fmt.Errorf("file is nill.")
// 	}

// 	path, _ := os.Getwd()
// 	// fileName := "example_local.yaml"
// 	fullPath := filepath.Join(path, fileName)

// 	fileBody, err := ioutil.ReadFile(fullPath)
// 	if err != nil {
// 		return fmt.Errorf("err:%s filename:%s", err.Error(), fileName)
// 	}

// 	var fileKey strings.Builder
// 	fileKey.WriteString(cfg.etcdCli.GetPrefix())
// 	fileKey.WriteString("/")
// 	fileKey.WriteString("aws-sg/dev/")
// 	fileKey.WriteString("repo/appsvc/v1")

// 	cfg.rootKey = fileKey.String()
// 	// fileKey := strings.Join(e.prefix, "/", "confd/aws-sg/dev/repo/appsvc/v1")
// 	key1, value1 := cfg.rootKey, string(fileBody)
// 	ctx, cancel := context.WithTimeout(context.Background(), cfg.etcdCli.GetRequestTimeout())
// 	defer cancel()
// 	if _, err := cfg.etcdCli.Client().Put(ctx, key1, value1); err != nil {
// 		return fmt.Errorf("err:%s key:%s", err.Error(), key1)
// 	}
// 	fmt.Println("fileKey:", cfg.rootKey)

// 	return nil
// }

// //loadConfigFromEtcd 依server的rootkey，取出对应的配置文件，并解析更新当前配置
// func (cfd *Confd) loadConfigFromEtcd() error {

// 	// cfg.mu.RLock()
// 	// defer cfg.mu.RUnlock()

// 	gresp, err := cfg.etcdCli.Client().Get(context.TODO(), cfg.rootKey, clientv3.WithPrefix())
// 	if err != nil {
// 		return err
// 	}

// 	var configData *bytes.Buffer
// 	for _, ev := range gresp.Kvs {
// 		configData = bytes.NewBuffer(ev.Value)
// 		fmt.Printf("[Kvs] %s : %s\n", ev.Key, ev.Value)
// 	}

// 	viper.SetConfigType("yaml")
// 	viper.ReadConfig(configData)

// 	xx := viper.Get("Service.Name")
// 	fmt.Println(" Service.Name", xx)
// 	fmt.Println(" namespace:", viper.GetString("Namespace"))

// 	return nil
// }

// func (cfd *Confd) loadConfigFromLocalFile() error {
// 	fmt.Println("[loadConfigFromLocalFile] ---")
// 	return nil
// }

// //WatchRemoteConfig 监控ETCD,保存配置热更新
// func (cfd *Confd) WatchRemoteConfig() {
// 	f1 := func() {
// 		for {
// 			rch := cfg.etcdCli.Client().Watch(context.Background(), cfg.rootKey, clientv3.WithPrefix())
// 			for wresp := range rch {
// 				for _, ev := range wresp.Events {
// 					// fmt.Printf("[Watch] %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
// 					switch ev.Type {
// 					case clientv3.EventTypePut:
// 						fmt.Printf("[Watch]:EventTypePut: [%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

// 					case clientv3.EventTypeDelete:
// 						fmt.Printf("[Watch]:EventTypeDelete: [%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

// 					}

// 				}
// 			}
// 		} //watch
// 	} // f1
// 	cfg.watchOnce.Do(f1)
// }

/*
	gresp, err := cfg.etcdCli.Client().Get(ctx, cfg.rootKey, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return err
	}
	//
	// fmt.Println("keys:", gresp.Kvs)
	for _, ev := range gresp.Kvs {
		fmt.Printf("[Kvs] %s : %s\n", ev.Key, ev.Value)
	}
*/
// func (cfd *Confd) Connect(cli *etcd.EtcdCli) error {

// 	cfg.etcdCli = etcd.NewEtcdCli(
// 		etcd.WithAddress([]string{"localhost:2379"}), //"localhost:2379"
// 		etcd.WithDialTimeout(2*time.Second),
// 		etcd.WithRequestTimeout(1*time.Second),
// 		etcd.WithPrefix("/confd/"))

// 	err := cfg.etcdCli.Connect()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (cfd *Confd) BindConfig() error {
// 	// c.confType = "ETCD" //"FILE"
// 	switch cfg.confType {
// 	case CETCD:
// 		return cfg.loadConfigFromEtcd()
// 	case CFILE:
// 		return cfg.loadConfigFromLocalFile()
// 	}
// 	return nil
// }

// func (cfd *Confd) SetConfigType(t string) error {
// 	switch t {
// 	case CETCD:
// 		cfg.confType = CETCD
// 	case CFILE:
// 		cfg.confType = CFILE
// 	default:
// 		return fmt.Errorf("Unknown : %s", t)
// 	}
// 	return nil
// }
