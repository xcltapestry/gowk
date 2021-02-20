package istio
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
/*
 使用场景说明:
    如果应用上了service mesh且跑在Istio 1.7以下环境，可以在application启动代码中，加上
istio.HoldApplicationUntilProxyStarts() 检查，以保证sidecar比应用先启动，否则可
能出现应用先启动但它的请求发不出去。
   否则，不需要使用它。
   
另一种情况：
    在停止服务时，sidecar比应用先停，这种情况也会出现应用的请求发不出去。这个代码解决不了。

最终的解决方案：
  在kubernets 1.18版本,k8s已考虑到sidecar，所以这两种情况都已解决。但实际上第一需要托管的各云服务商的k8s支持，第二要考虑升级风险。

*/


import (
	"net/http"
	"os"
	"time"

	"github.com/xcltapestry/gowk/pkg/logger"
)


// Delaying application start until sidecar is ready
// 基础库加代码仅适合 Istio 1.6
// Istio 1.7可通过配置 --set values.global.proxy.holdApplicationUntilProxyStarts=true
func HoldApplicationUntilProxyStarts() error {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	if len(host) <= 0 { // 非k8s环境
		return nil
	}

	if len(os.Getenv("NO_SIDECAR")) > 0 { // 没有启用Sidecar的，跳过检查
		return nil
	}

	var timeoutSeconds int = 30
	var periodMillis int = 10
	var url string = "http://localhost:15020/healthz/ready" // Envoy 的健康检查接口

	var err error
	var result bool
	timeoutAt := time.Now().Add(time.Duration(timeoutSeconds) * time.Second)
	for time.Now().Before(timeoutAt) {
		result, err = checkIfReady(url)
		if err == nil && result == true {
			logger.Infof("Envoy is ready!")
			return nil
		}
		logger.Debugf("Not ready yet: %v", err)
		time.Sleep(time.Duration(periodMillis) * time.Millisecond)
	}
	return nil
}

func checkIfReady(url string) (bool, error) {
	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	// Envoy 的健康检查通过，xDS 配置初始化完成
	if resp.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}

