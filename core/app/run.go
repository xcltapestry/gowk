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

import (
	"context"
	"runtime"

	"github.com/xcltapestry/gowk/core/confd"
	"github.com/xcltapestry/gowk/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

var App *Application

func New() *Application {

	//_ "go.uber.org/automaxprocs" Automatically set GOMAXPROCS to match Linux container CPU quota.
	_, _ = maxprocs.Set(maxprocs.Logger(logger.Infof))
	logger.Infof("runtime.NumCPU: %d runtime.GOMAXPROCS: %d  ", runtime.NumCPU(), runtime.GOMAXPROCS(-1))

	App = NewApplication()
	err := App.LoadConfig()
	if err != nil {
		logger.Fatalw("配置读取失败", " err:", err.Error())
	}
	logger.Infof("应用配置信息: %s", App.Confd().String())

	return App
}

func Serve(s ...Server) error {
	return App.Serve(s...)
}

func Run() error {
	return App.Run()
}

func Stop(ctx context.Context) {

}

func Flush() {
	App.Flush()
}

func Confd() *confd.Confd {
	return App.Confd()
}
