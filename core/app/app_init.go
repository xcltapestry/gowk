package app

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

// type ServiceDiscovery struct {
// }

import (
	"flag"

	"github.com/xcltapestry/gowk/core/confd"
)

func (app *Application) init() error {

	appCfg, err := app.parseFlags()
	if err != nil {
		return err
	}
	app.confdx = confd.NewConfd(appCfg)
	return nil
}

func (app *Application) parseFlags() (*confd.AppConfig, error) {

	var (
		deployEnv             string
		namespace             string
		appName               string
		version               string
		logsvcNats            string
		confdLocalFile        string
		confdRemoteAddrs      string
		confdRemoteConfigType string
		confdSyncMode         bool
	)

	flag.StringVar(&deployEnv, "deployenv", "dev", "dev/prod/pre/uat")
	//如果confd.remote.addrs为空或不通，则查找本地yaml.如都找不到，服务启动失败，如都有配置，优化找远程配置中心
	flag.StringVar(&confdLocalFile, "confd.local.file", "config.yaml", "本地配置文件如: config.yaml")
	flag.StringVar(&confdRemoteAddrs, "confd.remote.addrs", "", "远程配置中心地址，使用;分隔，如ETCD地址localhost:2379")
	flag.StringVar(&confdRemoteConfigType, "confd.remote.configtype", "yaml", "指定配置中心还原解析时的文件类型，默认yaml")
	flag.BoolVar(&confdSyncMode, "confd.remote.watch", false, "是否监控配置变更自动同步？默认为否.")
	flag.StringVar(&namespace, "namespace", "default", "")
	flag.StringVar(&appName, "app.name", "default", "")
	flag.StringVar(&version, "app.version", "01", "")
	flag.StringVar(&logsvcNats, "logsvc.nats", ";;nats://127.0.0.1:7222", "日志上报: nats标准配置,用;分隔")

	flag.Parse()

	scfg := confd.NewAppConfig()
	err := scfg.GetAppFlags(deployEnv, namespace, appName, version,
		confdLocalFile, confdRemoteAddrs, confdRemoteConfigType, confdSyncMode, logsvcNats)
	if err != nil {
		return nil, err
	}
	return scfg, nil
}

func (app *Application) LoadConfig() error {

	err := app.confdx.BindConfig()
	if err != nil {
		return err
	}

	// fmt.Println(" ------------------------------------- ")
	// keys := app.confdx.AllKeys()
	// for _, k := range keys {
	// 	fmt.Printf("AllKeys----[keys] %s\n", k)
	// }

	// fmt.Println("[confd] rootkey:", app.confdx.GetRootKey())
	// fmt.Println("[confd] Service.Name", app.confdx.Get("Service.Name"))
	// fmt.Println("[confd] namespace:", app.confdx.GetString("Namespace"))
	// // fmt.Println("[confd] Env:", app.Confd.GetString(confdx.CONFIG_env))
	// fmt.Println("[confd] service.redis.port:", app.confdx.GetString("service.redis.port"))

	return nil
}
