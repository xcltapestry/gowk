package apps

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
	"sync"

	"github.com/xcltapestry/gowk/pkg/config"
)

type Application struct {
	servers                         []Server
	initOnce, startupOnce, stopOnce sync.Once

	Config *config.Config
	Meta   *Metadata
}

func (app *Application) Serve(s ...Server) error {
	if app == nil {
		return fmt.Errorf("app is nill.")
	}
	app.servers = append(app.servers, s...)
	return nil
}

func (app *Application) Run() error {

	if app == nil {
		return fmt.Errorf("app is nill.")
	}

	if len(app.servers) <= 0 {
		os.Exit(0)
	}

	// app.buildConfig()

	app.startupOnce.Do(func() {
		var err error
		for _, s := range app.servers {
			svc := s
			//Todo, 并行运行和统一信号量关闭
			err = svc.Run()
			if err != nil {
				os.Exit(1)
			}
		}
	})

	return nil
}

func (app *Application) RegisterService() {

}

func (app *Application) UnregisterService() {

}

func (app *Application) buildConfig() {
	if app.Config == nil {
		return
	}
	if app.Meta == nil {
		app.Meta = NewMetadata()
	}

	if app.Config.InConfig("App.Namespace") {
		app.Meta.Namespace = app.Config.GetString("App.Namespace")
	} else {
		app.Meta.Namespace = "default"
	}

	if app.Config.InConfig("App.Id") {
		app.Meta.Id = app.Config.GetString("App.Id")
	} else {
		app.Meta.Id = "01"
	}

	if app.Config.InConfig("Env") {
		app.Meta.Env = app.Config.GetString("Env")
	} else {
		app.Meta.Env = "dev"
	}

}

func (app *Application) GetMetadata() *Metadata {
	if app.Meta == nil {
		app.Meta = NewMetadata()
	}
	return app.Meta
}

//
type ApplicationOption func(*Application)

func NewApp(options ...func(*Application)) *Application {
	app := &Application{}
	app.Config = config.New()
	app.Meta = NewMetadata()

	for _, f := range options {
		f(app)
	}

	return app
}

func InitConfig(path, fileName string) ApplicationOption {
	return func(a *Application) {
		a.Config.ReadInConfig(path, fileName)
		a.buildConfig()
	}
}

func Register(c string) ApplicationOption {
	return func(a *Application) {
		fmt.Println("[Register] c:", c)
		a.RegisterService()
	}
}
