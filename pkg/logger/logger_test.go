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

import (
	"errors"
	"flag"
	"testing"

	"github.com/golang/glog"
)

//  go test -v -cover=true

func TestLogger(t *testing.T) {

	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	WithName("notifysrv")
	WithValues("cloud", "aws-sg", "name", "notifysrv.xcl-servicestat:8080", "instanceIP", "172.168.1.1")

	Error(errors.New("Request failed with status code 503"), "handlecreate()", "requestId", 1, "parm", map[string]int{"k": 1})
	Info("name resolution failed", "val1", 1, "val2", map[string]int{"k": 1})

	/*
		output:
			notifysrv "handlecreate()" err="Request failed with status code 503" requestId=1 parm=map[k:1] cloud="aws-sg" name="notifysrv.xcl-servicestat:8080" instanceIP="172.168.1.1" error="Request failed with status code 503"
			notifysrv "name resolution failed" val1=1 val2=map[k:1] instanceIP="172.168.1.1" cloud="aws-sg" name="notifysrv.xcl-servicestat:8080"
	*/
}
