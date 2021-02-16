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
 */

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xcltapestry/gowk/core/naming"

	pb "github.com/xcltapestry/gowk/_example/protocol"
)

var _serviceName string = "HelloService"
var _etcdAddr string = "localhost:2379"

func main() {
	discovery := naming.NewNaming()
	err := discovery.GetResolver(_serviceName)
	if err != nil {
		log.Panic(err, " serviceName:", _serviceName)
	}

	for i := 0; i < 10; i++ {
		request := &pb.HelloRequest{Greeting: fmt.Sprintf("%d", i)}
		//从etcd获取rpc连接
		rpcConn, err := discovery.GetClientConn(_serviceName)
		if err != nil {
			log.Fatal(err, " _serviceName:", _serviceName)
		}
		client := pb.NewHelloServiceClient(rpcConn)
		//发送信息
		resp, err := client.SayHello(context.Background(), request)
		if err != nil {
			fmt.Println(" err:", err)
		} else {
			fmt.Println(" resp:", resp)
		}
		fmt.Printf("test => resp: %+v, err: %+v\n", resp, err)
		time.Sleep(time.Second)
	}

}
