# gowk
Service framework

#### 终于有精力和心情整下自己要用的服务框架

 - "少即是多"，简洁的框架才可以更方便的和各种基础设施集成
 - 技术框架和业务轮子的边界要清楚
 - 框架开放能力的"收"与"放"，需要多花时间思考
 - 框架可以提供些简洁好用的工具链，降低使用门槛
 - 基础设施(k8s? k8s+service mesh?)处于哪个阶段，一定程度上决定了框架的"厚薄"。
 - 云原生、k8s、Mesh玩法不同了，要融入并扩展这些基础能力到框架
 - 一开始就要考虑，如果一些当初认为"很棒"的设计被大量使用后，你想用"更棒"的设计来升级时如何做？
 - 如果有信心，可以在框架设计一开始就考虑，多云服务商、多数据中心、大数据的数据来源等看似比较"长远"的问题
 - 框架大部份情况下是很个性化的创造，要有独立思考的能力


#### Init module
  go mod init github.com/xcltapestry/gowk


#### 感谢
 - 表仅只列出了主要的，有直接使用或有所借鉴的开源项目，感谢所有参与开源的人们
  
| 分类 |  url | 备注 |
| :---- | :---- | :---- | 
| redis | github.com/go-redis | |
| redis | github.com/bsm/redislock |  |
| log |  github.com/golang/glog |  |
| log |  github.com/kubernetes/klog | 参考 |
| log |  github.com/go-logr/logr | 参考 |
| log | github.com/golang/glog |   |

