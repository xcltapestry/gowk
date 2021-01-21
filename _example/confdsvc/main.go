package main

import (

	"fmt"
	"time"

	
	// "github.com/xcltapestry/gowk/core/confd"


	// "github.com/xcltapestry/gowk/core/apps"

	// "/Users/xcl/Documents/myworkspace/frameworks/gowk/core/apps"

	apps "../../core/apps"

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

	glog.Info("golog ------------- ")

	fmt.Println(" ok ")

	//apps.New().Run()

	defer apps.Flush()

	time.Sleep(1 * time.Second)
}
