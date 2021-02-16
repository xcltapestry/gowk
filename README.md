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
 - 框架大部份情况下是很个性化的工程性创造，要有独立思考的能力


#### Init module
  go mod init github.com/xcltapestry/gowk

#### 例子


服务端（_example/rpc_http/service/service.go）:
```go

 func main() {
	logger.NewDefaultLogger()

	app.New()
	
	// http 服务
	httpSvc := services.NewHTTPService()
	httpSvc.Router(RegisterHandlers)
	httpSvc.Listen(app.Confd().GetString("application.listen.http"))
	app.Serve(httpSvc)

	// gRPC 服务
	rpcSvc := CreateRPC()
	// 注册 hello RPC服务
	RegisterHelloServer(rpcSvc)
	app.Serve(rpcSvc)

	app.Run()
}

func CreateRPC() *services.RPCService {
	// rpc服务
	rpcSvc := services.NewRPCService()
	_, err := rpcSvc.Listen()
	if err != nil {
		logger.Panic(err)
	}
	// 从配置中心获取配置
	etcdAddrs := app.Confd().GetStringSlice("application.registry.etcdv3.endpoints")
	etcdDialTimeout := app.Confd().GetDuration("application.registry.etcdv3.dialtimeout")
	// 基于etcd的RPC名字服务
	rpcSvc.NewNaming(naming.WithAddress(etcdAddrs), naming.WithDialTimeout(etcdDialTimeout))
	return rpcSvc
}

func RegisterHelloServer(rpcSvc *services.RPCService) {
	// hello 服务
	if err := rpcSvc.Registry("HelloService"); err != nil {
		logger.Panic(err)
	}
	// 注册Hello服务
	pb.RegisterHelloServiceServer(rpcSvc.GetgRPCServer(), &HelloServer{})
}

//RegisterHandlers 路由
func RegisterHandlers(m *mux.Router) {
	m.HandleFunc("/health", HealthHandler)
	m.Handle("/", http.NotFoundHandler())
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug(" Health => ", time.Now().Unix())
	fmt.Fprintf(w, "Health: %v\n", time.Now().Unix())
}

//HelloServer type
type HelloServer struct{}

//SayHello func
func (h *HelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	response := &pb.HelloResponse{
		Reply: fmt.Sprintf("hello, %s", req.Greeting),
	}
	fmt.Println(req.Greeting)
	return response, nil
}

```
配置文件：
```yaml
ApiVersion: 0.0.1
Env: dev
Namespace: xxx.cluster
application:
  Listen:
      http: :8088
      grpc: :9099
  Registry:
      etcdv3:
        endpoints: ["127.0.0.1:2379"]
        dialTimeout: 10s
Service:
  Name: ordersvc
  Redis:
    Addr: 127.0.0.1
    Port: :6377
  

```

#### 感谢
 - 表仅只列出了主要的，有直接使用或有所借鉴的开源项目，感谢所有参与开源的人们
  
| 分类 |  url | 备注 |
| :---- | :---- | :---- | 
| redis | github.com/go-redis | |
| redis | github.com/bsm/redislock |  |
| log | go.uber.org/zap |  |
| log | github.com/natefinch/lumberjack |  |
| confd | go.etcd.io/etcd/clientv3 |  |
| confd | github.com/spf13/viper |   |
| router | github.com/gorilla/mux |   |
| json | github.com/json-iterator/go |   |






