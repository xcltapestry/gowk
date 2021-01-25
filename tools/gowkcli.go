package main 

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
 
/*
////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//   配置中心
//
////////////////////////////////////////////////////////////////////////////////////////////////////////
远程配置中心处理流程：
   1. 读取配置文件信息，存储到远程配置中心指定Key上
   2. 服务运行时，依所传的flags参数，得到Key后，从远程配置中心获取内容
   3. viper 解析内容，并映射到每个配置项


服务使用备注:
    服务配置参数:
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

    服务Key:
    _ConfDPrefix = "confd"
	_ConfigType  = "yaml"
	c.rootKey = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s/", _ConfDPrefix, c.cloud, c.region,
		c.deployEnv, c.namespace, c.appName, c.appVersion)

////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//   命令行工具
//
////////////////////////////////////////////////////////////////////////////////////////////////////////
go run main.go --help

命令行参数:

flag:
  -confd.remote.addrs string
    	远程配置中心地址，使用;分隔，如ETCD地址localhost:2379 (default "localhost:2379")


案例:

Use the arrow keys to navigate: ↓ ↑ → ←
? 选择::
  ▸ 1.上转指定配置文件至ETCD
    2.查看所有配置


✔ 1.上转指定配置文件至ETCD
You choose "1.上转指定配置文件至ETCD "
输入配置文件全路径: /Users/xxx/Documents/workspace/mytest/etcdt2cmd/example_local.yaml
输入需要查看的Etcd Key: aa
确认执行(y/n): y
Key已存在，是否替换(y/n)?: y
配置已更新,从配置文件同步到Etcd成功！
----------------------------
ETCD Key:  aa


✔ 2.查看所有配置
You choose "2.查看所有配置"
输入需要查看的Etcd Key: aa
确认执行(y/n): y
从ETCD取值成功！
......
////////////////////////////////////////////////////////////////////////////////////////////////////////

*/

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
)

var confdRemoteAddrs string

func main() {
	flag.StringVar(&confdRemoteAddrs, "confd.remote.addrs", "localhost:2379", "远程配置中心地址，使用;分隔，如ETCD地址localhost:2379")
	flag.Parse()
	cmdsMain()
}

func cmdsMain() {

CHOOSE:
	prompt := promptui.Select{
		Label: "选择:",
		Items: []string{"1.上转指定配置文件至ETCD ", "2.查看所有配置"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("You choose %q\n", result)

	switch {
	case strings.Contains(result, "1."):
		cmdsStep1Ask()
	case strings.Contains(result, "2."):
		cmdsStep2Ask()
	default:
		goto CHOOSE
	}
}

//cmdsStep1Ask 1.上转指定配置文件至ETCD
func cmdsStep1Ask() {
	respConfFile := cmdsInput("输入配置文件全路径: ")
	respEtcdKey := cmdsInput("输入需要查看的Etcd Key: ")

	var respYN string
YN:
	respYN = cmdsInput("确认执行(y/n): ")
	check, exec := confirmation(respYN)
	if !check {
		goto YN
	}
	if !exec {
		return
	}
	execChoose1(respConfFile, respEtcdKey)
}

//cmdsStep2Ask 2.查看所有配置
func cmdsStep2Ask() {
	respEtcdKey := cmdsInput("输入需要查看的Etcd Key: ")
	var respYN string
YN:
	respYN = cmdsInput("确认执行(y/n): ")
	check, exec := confirmation(respYN)
	if !check {
		goto YN
	}
	if !exec {
		return
	}
	execChoose2(respEtcdKey)
}

//execChoose1 将本地配置文件，推到远程ETCD指定Key下
func execChoose1(respConfFile, respEtcdKey string) {
	// 检查文件是否存在
	_, _, err := PaseConfigFile(respConfFile)
	if err != nil {
		fmt.Println("发生异常！ err: ", err)
		return
	}

	fileBody, err := ioutil.ReadFile(respConfFile)
	if err != nil {
		fmt.Println("err:%s confFile:%s", err.Error(), respConfFile)
		return
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{confdRemoteAddrs},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		fmt.Println("发生异常！ err: ", err)
		return
	}
	defer cli.Close()

	// 执行
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel1()
	gresp, err := cli.Get(ctx1, respEtcdKey, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("发生异常！ err: ", err)
		return
	}
	if gresp.Count > 0 { // 检查Key是否已存在，如存在，需选择是否要替换?
	YN:
		respYN := cmdsInput("Key已存在，是否替换(y/n)?: ")
		check, exec := confirmation(respYN)
		if !check {
			goto YN
		}
		if !exec {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if _, err := cli.Put(ctx, respEtcdKey, string(fileBody)); err != nil {
		fmt.Println("err:%s rootKey:%s", err.Error(), respEtcdKey)
		return
	}

	fmt.Println("配置已更新,从配置文件同步到Etcd成功！ ")
	fmt.Println("----------------------------")
	fmt.Println("ETCD Key: ", respEtcdKey)
}

//execChoose2 查看远程Key的配置
func execChoose2(respEtcdKey string) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{confdRemoteAddrs},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		fmt.Println("发生异常！ err: ", err)
		return
	}
	defer cli.Close()

	// 执行
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	gresp, err := cli.Get(ctx, respEtcdKey, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("发生异常！ err: ", err)
		return
	}
	if gresp.Count == 0 {
		fmt.Println(" 在ETCD中没有找到数据！key:", respEtcdKey)
		return
	}

	fmt.Println("从ETCD取值成功！ ")
	fmt.Println("----------------------------")
	for _, item := range gresp.Kvs {
		fmt.Printf("%s : %s \n", item.Key, item.Value)
	}
	fmt.Println("----------------------------")
	fmt.Println("ETCD Key: ", respEtcdKey)

	// loadViper := viper.New()
	// configType := "yaml"
	// loadViper.SetConfigType(configType)

	// var configData *bytes.Buffer
	// configData = bytes.NewBuffer(gresp.Kvs[0].Value)
	// err = loadViper.ReadConfig(configData)
	// if err != nil {
	// 	fmt.Println("发生异常！ err: ", err)
	// 	return
	// }

	// fmt.Println("对应的conf配置Key如下: ")
	// confs := loadViper.AllKeys()
	// for i, v := range confs {
	// 	fmt.Printf("%d : %s \r\n", i, v)
	// }
}

func cmdsInput(displayText string) string {
	var resp string
	fmt.Print(displayText)
	_, err := fmt.Scanln(&resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func confirmation(respYN string) (check bool, exec bool) {
	switch strings.ToLower(respYN) {
	case "y", "yes": //执行下一步"
		check, exec = true, true
	case "n", "no":
		check, exec = true, false
	default:
		check, exec = false, false
	}
	return
}

func PaseConfigFile(confFile string) (string, string, error) {

	s, err := os.Stat(confFile)
	if err != nil && os.IsNotExist(err) {
		return "", "", fmt.Errorf("需输入确实存在的文件全路径！ err:%s", err)
	}
	if s.IsDir() {
		return "", "", fmt.Errorf("需输入文件全路径，不能输入目录！ 目录:%s", confFile)
	}
	path, fileName := filepath.Split(confFile)
	ext := strings.ToLower(filepath.Ext(fileName))
	fext := SubString(ext, 1, len(ext))

	return path, fext, nil
}

//包含中文的字符串截取
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}