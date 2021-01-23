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
	"fmt"
	"log"

	"github.com/xcltapestry/gowk/core/confd"
)

var App *Application

func New() *Application {

	App = NewApplication()
	err := App.LoadConfig()
	if err != nil {
		fmt.Println(" err:", err)
		log.Fatal(" err:", err)
	}
	fmt.Println(App.Confd().String())
	return App
}

func Run() error {
	return App.Run()
}

func Flush() {
	App.Flush()
}

func Confd() *confd.Confd {
	return App.Confd()
}
