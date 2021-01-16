package waitgroup

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
	"sync"
)

//WaitGroupWrapper 并行管理器
type WaitGroupWrapper struct {
	sync.WaitGroup
}

//Wrap 执行函数
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

//WaitFunc 函数管理
func WaitFunc(arr []func()) {
	var wg WaitGroupWrapper
	for _, v := range arr {
		wg.Wrap(v)
	}
	wg.Wait()
}

// WaitGroupFunc 支持带有参数的函数
type WaitGroupFunc struct {
	Func func(args interface{})
	Args interface{}
}

//WrapWithArgs 执行函数, 支持传参的协程
func (w *WaitGroupWrapper) WrapWithArgs(cb func(args interface{}), argsItem interface{}) {
	w.Add(1)
	go func(args interface{}) {
		cb(args)
		w.Done()
	}(argsItem)
}

//WaitFuncWithArgs 函数管理, 支持传参的协程
func WaitFuncWithArgs(arr []WaitGroupFunc) {
	var wg WaitGroupWrapper
	for _, v := range arr {
		wg.WrapWithArgs(v.Func, v.Args)
	}
	wg.Wait()
}
