package confd



import (
	"fmt"
	"flag"
	"testing"
)


// go test -v -run="TestFlags" -args -deployenv=true
func TestFlags(t *testing.T) {
	//testing.Init()

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
	flag.StringVar(&confdLocalFile, "confd.local.file", "conf_test.yaml", "本地配置文件如: config.yaml")
	flag.StringVar(&confdRemoteAddrs, "confd.remote.addrs", "", "远程配置中心地址，使用;分隔，如ETCD地址localhost:2379")
	flag.StringVar(&confdRemoteConfigType, "confd.remote.configtype", "yaml", "指定配置中心还原解析时的文件类型，默认yaml")
	flag.BoolVar(&confdSyncMode, "confd.remote.watch", false, "是否监控配置变更自动同步？默认为否.")
	flag.StringVar(&namespace, "namespace", "default", "")
	flag.StringVar(&appName, "app.name", "default", "")
	flag.StringVar(&version, "app.version", "01", "")
	flag.StringVar(&logsvcNats, "logsvc.nats", ";;nats://127.0.0.1:7222", "日志上报: nats标准配置,用;分隔")
	flag.Parse()


// conf_test.yaml

	fmt.Println("deployEnv:", deployEnv)
	fmt.Println("namespace:", namespace)
	fmt.Println("appName:", appName)
	fmt.Println("version:", version)
	

	scfg := NewAppConfig()
	scfg.UpdateConfig(deployEnv, namespace, appName, version, confdLocalFile, confdRemoteAddrs, confdRemoteConfigType, confdSyncMode, logsvcNats)

	confdx := NewConfd(scfg)

	err2 := confdx.ReadConfigFileToRemote()
	fmt.Println("ReadConfigFileToRemote  err:", err2)
	err := confdx.BindConfig()
	fmt.Println("BindConfig  err:", err)

	fmt.Println(" ------------------------------------- ")
	keys := confdx.AllKeys()
	for _, k := range keys {
		fmt.Printf("AllKeys----[keys] %s\n", k)
	}




}