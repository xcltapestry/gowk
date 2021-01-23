package logger

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
	"bytes"
	"fmt"

	"github.com/golang/glog"
)

const missingValue = "(MISSING)"
const logDepth = 2

//glogx : 实现参考go-logr设计
type glogx struct {
	name      string
	keyValues map[string]interface{}
}

func (l *glogx) WithName(name string) *glogx {
	l.name = name
	return l
}

func (l *glogx) WithValues(kvs ...interface{}) *glogx {
	newMap := make(map[string]interface{}, len(l.keyValues)+len(kvs)/2)
	for k, v := range l.keyValues {
		newMap[k] = v
	}
	for i := 0; i < len(kvs); i += 2 {
		switch kvs[i].(type) {
		case string:
			newMap[kvs[i].(string)] = kvs[i+1]
		default:
			continue
		}
	}
	l.keyValues = newMap
	return l
}

//kvListFormat : 主要代码引用自k8s中看到的klog源码printS,目的是为了更好的结构化日志
func (l *glogx) printS(s severity, err error, msg string, keysAndValues ...interface{}) {
	b := &bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%q", msg))
	if err != nil {
		b.WriteByte(' ')
		b.WriteString(fmt.Sprintf("err=%q", err.Error()))
	}
	l.kvListFormat(b, keysAndValues...)

	// 使用glog打印具体日志
	n := fmt.Sprintf("%s ", l.name)
	switch s {
	case infoLog:
		glog.InfoDepth(logDepth, n, b)
	case warningLog:
		glog.WarningDepth(logDepth, n, b)
	case errorLog:
		glog.ErrorDepth(logDepth, n, b)
	case fatalLog:
		glog.FatalDepth(logDepth, n, b)
	default:
		glog.InfoDepth(logDepth, n, b)
	}
}

//kvListFormat : 代码引用自k8s所用中，klog源码的kvListFormat
func (l *glogx) kvListFormat(b *bytes.Buffer, keysAndValues ...interface{}) {
	for i := 0; i < len(keysAndValues); i += 2 {
		var v interface{}
		k := keysAndValues[i]
		if i+1 < len(keysAndValues) {
			v = keysAndValues[i+1]
		} else {
			v = missingValue
		}
		b.WriteByte(' ')

		switch v.(type) {
		case string, error:
			b.WriteString(fmt.Sprintf("%s=%q", k, v))
		default:
			if _, ok := v.(fmt.Stringer); ok {
				b.WriteString(fmt.Sprintf("%s=%q", k, v))
			} else {
				b.WriteString(fmt.Sprintf("%s=%+v", k, v))
			}
		}
	}
}
