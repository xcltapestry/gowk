package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/xcltapestry/gowk/core/app"
	"github.com/xcltapestry/gowk/core/services"
	"github.com/xcltapestry/gowk/pkg/logger"
)

var (
	// go build 时通过 ldflags 设置
	codeBuiildSourceVersion = ""
	codeBuildTime           = ""
	gitHash                 = ""
)

/*
test:
	go run main.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -confd.remote.addrs="localhost:2379"

	go run main.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -confd.remote.addrs=""

	go run main.go -deployenv=prod  -namespace=order -app.name=ordersvc -app.version=02 -confd.local.file=conf.yaml -logger.output=alsologtostdout

*/
func main() {

	// logger.NewDefaultLogger()
	logger.NewLogger(logger.Text)
	app.New()
	httpsvc := services.NewHTTPService()
	httpsvc.Router(RegisterHandlers)
	httpsvc.SetHTTPAddr(":8003")
	app.Serve(httpsvc)
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
