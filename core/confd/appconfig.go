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
	"os"
	"strings"
)

type AppConfig struct {
	// -- env --
	cloud  string
	region string

	// -- build --
	codeBuiildSourceVersion string
	codeBuildTime           string
	gitHash                 string

	// -- flag  - application
	deployEnv  string
	namespace  string
	appName    string
	appVersion string
	logsvcNats string // 日志上报服务器地址

	// --- flag - confd
	// confd.remote.addrs与localConfigFile 必填一个，否则启动失败
	// 如果confd.remote.addrs为空或不通，则通过localConfigFile查找本地yaml.
	localConfd            bool
	localConfigFile       string
	confdRemoteAddrs      string
	confdRemoteConfigType string
	confdSyncMode         bool // 默认不打开与远端的自动同步

	// --- etcd key
	rootKey string
}

const (
	_cloud  = "default" //"aws"  aws ec2 describe-regions
	_region = "default" //"eu-north-1

	_envCloud  = "X_APP_CLOUD"
	_envRegion = "X_APP_REGION"

	_ConfDPrefix = "confd"
	_ConfigType  = "yaml"
)

func NewAppConfig() *AppConfig {
	return &AppConfig{cloud: _cloud, region: _region,
		localConfd: true,
		deployEnv:  "undefined", namespace: "default", appName: "undefined", appVersion: "undefined"}
}

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		cloud:                 _cloud,
		region:                _region,
		localConfd:            true,
		deployEnv:             "undefined",
		namespace:             "undefined",
		appName:               "undefined",
		appVersion:            "undefined",
		localConfigFile:       "config.yaml",
		confdRemoteAddrs:      "localhost:2379",
		confdSyncMode:         false,
		confdRemoteConfigType: _ConfigType,
	}
}

func (c *AppConfig) GetRootKey() string {
	if c.rootKey == "" {
		c.updateRootKey()
	}
	return c.rootKey
}

func (c *AppConfig) updateRootKey() string {
	c.rootKey = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s/", _ConfDPrefix, c.cloud, c.region,
		c.deployEnv, c.namespace, c.appName, c.appVersion)
	return c.rootKey
}

func (c *AppConfig) GetRemoteConfigType() string {
	if c.confdRemoteConfigType == "" {
		return _ConfigType
	}
	return c.confdRemoteConfigType
}

func (c *AppConfig) updateEnv() {
	appCloud := os.Getenv(_envCloud)
	appRegion := os.Getenv(_envRegion)
	if appCloud == "" {
		appCloud = _cloud
	}
	if appRegion == "" {
		appRegion = _region
	}
	c.cloud = appCloud
	c.region = appRegion
}

func (c *AppConfig) verify() error {
	if c.localConfigFile == "" && c.confdRemoteAddrs == "" {
		return fmt.Errorf("没有指定远程配置中心和本地配置文件。")
	}
	return nil
}

//UpdateConfig 从命令行参数与环境变量中获取参数
func (c *AppConfig) UpdateConfig(deployEnv, namespace, appName, appVersion string,
	confdLocalFile, confdRemoteAddrs, confdRemoteConfigType string, confdSyncMode bool,
	logsvcNats string) error {

	c.deployEnv, c.namespace, c.appName, c.appVersion = deployEnv, namespace, appName, appVersion
	c.localConfigFile, c.confdRemoteAddrs = confdLocalFile, confdRemoteAddrs
	c.confdRemoteConfigType = confdRemoteConfigType
	c.confdSyncMode = confdSyncMode
	c.logsvcNats = logsvcNats

	if err := c.verify(); err != nil {
		return err
	}

	if c.confdRemoteAddrs != "" {
		c.localConfd = false
	}

	if c.deployEnv == "" {
		c.deployEnv = "undefined"
	}

	if c.namespace == "" {
		c.namespace = "undefined"
	}

	if c.appName == "" {
		c.appName = "undefined"
	}

	if c.appVersion == "" {
		c.appVersion = "undefined"
	}

	if c.confdRemoteConfigType == "" {
		c.confdRemoteConfigType = _ConfigType
	}

	c.updateRootKey()
	c.updateEnv()

	return nil
}

func (c *AppConfig) IsLocal() bool {
	return c.localConfd
}

func (c *AppConfig) IsSync() bool {
	return c.confdSyncMode
}

func (c *AppConfig) GetLocalConfigFile() string {
	return c.localConfigFile
}

func (c *AppConfig) GetConfdRemoteAddrs() []string {
	//strings.Split("localhost:2379;localhost:2377;localhost:2376", ";")
	return strings.Split(c.confdRemoteAddrs, ";")
}

func (c *AppConfig) String() string {
	var builder strings.Builder
	builder.WriteString(" cloud=")
	builder.WriteString(c.cloud)
	builder.WriteString(" region=")
	builder.WriteString(c.region)
	builder.WriteString(" codeBuiildSourceVersion=")
	builder.WriteString(c.codeBuiildSourceVersion)
	builder.WriteString(" codeBuildTime=")
	builder.WriteString(c.codeBuildTime)
	builder.WriteString(" gitHash=")
	builder.WriteString(c.gitHash)
	builder.WriteString(" deployEnv=")
	builder.WriteString(c.deployEnv)
	builder.WriteString(" namespace=")
	builder.WriteString(c.namespace)
	builder.WriteString(" appName=")
	builder.WriteString(c.appName)
	builder.WriteString(" appVersion=")
	builder.WriteString(c.appVersion)
	builder.WriteString(" logsvcNats=")
	builder.WriteString(c.logsvcNats)
	builder.WriteString(" localConfigFile=")
	builder.WriteString(c.localConfigFile)
	builder.WriteString(" confdRemoteAddrs=")
	builder.WriteString(c.confdRemoteAddrs)
	builder.WriteString(" confdRemoteConfigType=")
	builder.WriteString(c.confdRemoteConfigType)
	builder.WriteString(" rootKey=")
	builder.WriteString(c.rootKey)
	return builder.String()
}
