package logger

/**
 * Copyright 2021  gowk Author. All Rights Reserved.
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

/*
ToDo:
  - 日志不落盘支持
*/

type severity int32

const (
	infoLog severity = iota
	warningLog
	errorLog
	fatalLog
	numSeverity = 4
)

var gLogger *glogx

func init() {
	gLogger = &glogx{}
}

//WithName -
func WithName(name string) {
	gLogger.WithName(name)
}

//WithValues -
func WithValues(kvs ...interface{}) {
	gLogger.WithValues(kvs...)
}

//Info 级别日志
func Info(msg string, kvs ...interface{}) {
	if len(gLogger.keyValues) > 0 {
		for k, v := range gLogger.keyValues {
			kvs = append(kvs, k, v)
		}
	}
	gLogger.printS(infoLog, nil, msg, kvs...)
}

//Error 级别日志
func Error(err error, msg string, kvs ...interface{}) {
	if len(gLogger.keyValues) > 0 {
		for k, v := range gLogger.keyValues {
			kvs = append(kvs, k, v)
		}
	}
	kvs = append(kvs, "error", err)
	gLogger.printS(errorLog, err, msg, kvs...)
}

//Warn 级别日志
func Warn(msg string, kvs ...interface{}) {
	if len(gLogger.keyValues) > 0 {
		for k, v := range gLogger.keyValues {
			kvs = append(kvs, k, v)
		}
	}
	gLogger.printS(infoLog, nil, msg, kvs...)
}

//Fatal 级别日志
func Fatal(err error, msg string, kvs ...interface{}) {
	if len(gLogger.keyValues) > 0 {
		for k, v := range gLogger.keyValues {
			kvs = append(kvs, k, v)
		}
	}
	kvs = append(kvs, "error", err)
	gLogger.printS(fatalLog, err, msg, kvs...)
}
