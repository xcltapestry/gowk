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
 * @Project gowk
 * @Description go framework
 * @author XiongChuanLiang<br/>(xcl_168@aliyun.com)
 * @license http://www.apache.org/licenses/  Apache v2 License
 * @version 1.0
 */

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/xcltapestry/gowk/core/apps"
	"github.com/xcltapestry/gowk/pkg/logger"
)

func main() {

	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	path, _ := os.Getwd()
	fileName := "example_local.yaml"

	var app *apps.Application
	//app = apps.New()
	//app.Config.ReadInConfig(path, fileName)

	app = apps.NewApp(
		apps.InitConfig(path, fileName))
	logger.WithName(app.Config.GetString("Service.WebApp.Name"))

	// httpsvc := apps.NewHTTPService()
	// httpsvc.Router(RegisterHandlers)
	// httpPort := app.Config.GetString("Service.WebApp.Port")
	// httpsvc.SetHTTPAddr(httpPort)
	// logger.Info("httpsvc ok httpPort:", httpPort)

	// app.Serve(httpSvc())

	httpsvc2 := apps.NewHTTPService()
	httpsvc2.Router(RegisterHandlers2)
	httpsvc2.SetHTTPAddr(":8003")
	app.Serve(httpsvc2)

	app.Run()
	logger.Info("end.")

}

func httpSvc() *apps.HTTPService {
	httpsvc := apps.NewHTTPService()
	httpsvc.Router(RegisterHandlers)
	// httpPort := app.Config.GetString("Service.WebApp.Port")
	// httpsvc.SetHTTPAddr(httpPort)
	return httpsvc
}

//RegisterHandlers 路由
func RegisterHandlers(m *mux.Router) {
	m.HandleFunc("/v1/ping", ArticlesCategoryHandler)
	m.HandleFunc("/v1/health", HealthHandler)
	m.Handle("/", http.NotFoundHandler())
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health: %v\n", time.Now().Unix())
}

//RegisterHandlers 路由
func RegisterHandlers2(m *mux.Router) {
	m.HandleFunc("/v1/ping", ArticlesCategoryHandler)
	m.HandleFunc("/v1/health", HealthHandler)
	m.Handle("/", http.NotFoundHandler())
}

// //ServiceConfig 服务本身的一些业务相关配置
// type ServiceConfig struct {
// 	Addr string
// }

// func (cfg *ServiceConfig) ReadConfig() error {

// 	return nil
// }
