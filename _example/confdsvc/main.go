package main

import (
	"flag"
	"fmt"
	"time"

	//confdx "github.com/xcltapestry/gowk/core/confd"

	"github.com/xcltapestry/gowk/core/confd"

	_ "go.uber.org/automaxprocs" //Automatically set GOMAXPROCS to match Linux container CPU quota.
)

var (
	// go build 时通过 ldflags 设置
	codeBuiildSourceVersion = ""
	codeBuildTime           = ""
	gitHash                 = ""
)

/*


go run main.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=example_local.yaml -confd.remote.addrs="localhost:2379"



*/
func main() {

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
	flag.StringVar(&confdLocalFile, "confd.local.file", "", "本地配置文件如: config.yaml")
	flag.StringVar(&confdRemoteAddrs, "confd.remote.addrs", "", "远程配置中心地址，使用;分隔，如ETCD地址localhost:2379")
	flag.StringVar(&confdRemoteConfigType, "confd.remote.configtype", "yaml", "指定配置中心还原解析时的文件类型，默认yaml")
	flag.BoolVar(&confdSyncMode, "confd.remote.watch", false, "是否监控配置变更自动同步？默认为否.")
	flag.StringVar(&namespace, "namespace", "default", "")
	flag.StringVar(&appName, "app.name", "default", "")
	flag.StringVar(&version, "app.version", "01", "")
	flag.StringVar(&logsvcNats, "logsvc.nats", ";;nats://127.0.0.1:7222", "日志上报: nats标准配置,用;分隔")
	flag.Parse()

	confd.NewAppConfig()
	fmt.Println(" ------------------------------------- ")


	/*
	scfg := confdx.NewAppConfig()
	scfg.UpdateConfig(deployEnv, namespace, appName, version, confdLocalFile, confdRemoteAddrs, confdRemoteConfigType, confdSyncMode, logsvcNats)

	confd := confdx.NewConfd(scfg)

	go func() {
		confd.WatchRemoteConfig()
	}()

	err2 := confd.ReadConfigFileToRemote()
	fmt.Println("ReadConfigFileToRemote  err:", err2)

	err := confd.BindConfig()
	fmt.Println("BindConfig  err:", err)

	fmt.Println(" ------------------------------------- ")
	keys := confd.AllKeys()
	for _, k := range keys {
		fmt.Printf("AllKeys----[keys] %s\n", k)
	}

	fmt.Println("[confd] rootkey:", confd.GetRootKey())
	fmt.Println("[confd] Service.Name", confd.Get("Service.Name"))
	fmt.Println("[confd] namespace:", confd.GetString("Namespace"))
	//	fmt.Println("[confd] Env:", confd.GetString(confdx.CONFIG_env))
	fmt.Println("[confd] service.redis.port:", confd.GetString("service.redis.port"))

	fmt.Println(" ------------------------------------- ")
*/
	time.Sleep(3 * time.Second)
}
