package main

import (
	"fmt"
	"net/http"
	"time"
	"context"

	"github.com/gorilla/mux"
	"github.com/xcltapestry/gowk/core/app"
	"github.com/xcltapestry/gowk/core/services"
	"github.com/xcltapestry/gowk/pkg/logger"

	pb "github.com/xcltapestry/gowk/_example/protocol"
)

var (
	// go build 时通过 ldflags 设置
	codeBuiildSourceVersion = ""
	codeBuildTime           = ""
	gitHash                 = ""
)

/*
test:
	go run service.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -confd.remote.addrs="localhost:2379"

	go run service.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -confd.remote.addrs=""

	go run main.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -logger.output=alsologtostdout

*/

func main() {
	logger.NewDefaultLogger()
	app.New()

	fmt.Println(" ------------------------------------- ")
	keys := app.Confd().AllKeys()
	for _, k := range keys {
		fmt.Printf("AllKeys----[keys] %s\n", k)
	}

	httpSvc := services.NewHTTPService()
	httpSvc.Router(RegisterHandlers)
	httpSvc.Listen(":8003")
	app.Serve(httpSvc)

	rpcSvc := services.NewRPCService()
	// 启动监听,初始化
	// gRPCServer,err := rpcSvc.Listen(":8082")
	gRPCServer,err := rpcSvc.Listen()
	if err != nil {
		logger.Infow("Serve --> err:",err )
		return 
	}
	// 注册服务
	pb.RegisterHelloServiceServer(gRPCServer, &HelloServer{})
	app.Serve(rpcSvc)

	logger.Infow("Serve -->Listen :8003")
	app.Run()
	logger.Infow("Serve --> end.")
}

//RegisterHandlers 路由
func RegisterHandlers(m *mux.Router) {
	m.HandleFunc("/ping", PingHandler)
	m.HandleFunc("/health", HealthHandler)
	m.Handle("/", http.NotFoundHandler())
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger.Debug(" Ping ---------- ")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ping: %v\n", vars["hello"])
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug(" Health ---------- ")
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

