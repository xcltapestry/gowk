package naming

/**

            +-------------+
            | etcd server |
            +------+------+
                   ^ watch key A (s-watcher)
                   |
           +-------+-----+
           | gRPC proxy  | <-------+
           |             |         |
           ++-----+------+         |watch key A (c-watcher)
watch key A ^     ^ watch key A    |
(c-watcher) |     | (c-watcher)    |
    +-------+-+  ++--------+  +----+----+
    |  client |  |  client |  |  client |
    |         |  |         |  |         |
    +---------+  +---------+  +---------+

https://etcd.io/docs/v3.1.12/op-guide/grpc_proxy/

**/


import (
	"fmt"
	"time"
	"net"

	"go.uber.org/zap"

	etcdClientV3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/proxy/grpcproxy"
)

type etcdRegister struct {
	ttl int
	cli         *etcdClientV3.Client
	prefix string 
}

func NewEtcdRegister()*etcdRegister{
	e := &etcdRegister{}
	e.ttl = 5
	e.cli = nil 
	e.prefix = "/discovery/etcd/grpc/service/"
	return e
}

func(r *etcdRegister)Register(serviceName,addr string)error {
	// 要保持 etcd 长连接 
	var err error 
	// config := etcdClientV3.Config{
	// 	Endpoints:   []string{"localhost:2379"},  //"localhost:2379"
	// 	DialTimeout: 3 * time.Second,
	// }

	// if r.cli == nil {
	// 	r.cli, err = etcdClientV3.New(config)
	// 	if err != nil {
	// 		return fmt.Errorf(" err:%s",err)
	// 	}
	// }

	r.cli, err = etcdClientV3.New(etcdClientV3.Config{
		Endpoints:   []string{"localhost:2379"},  //strings.Split(etcdAddr, ";"),
		DialTimeout: 15 * time.Second,
		LogConfig: &zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.ErrorLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:      "json",
			EncoderConfig: zap.NewProductionEncoderConfig(),
			// Use "/dev/null" to discard all
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	})
	if err != nil {
		return fmt.Errorf(" err:%s",err)
	}


	registerURL := r.GetRegisterURL(serviceName)
	grpcproxy.Register(r.cli,registerURL , addr, r.ttl)

	return nil 
}

func(r *etcdRegister)CloseRegisterEtcd(){
	if r.cli != nil {
		r.cli.Close()
	}
}

func(r *etcdRegister)GetRegisterURL(serviceName string)string{
	// "/discovery/etcd/grpc/service/%s"
	return fmt.Sprintf("%s%s",r.prefix  ,serviceName)
}


type Server struct {
	IP   net.IP
	Port int32
}

func (s Server) Equal(x Server) bool {
	return s.Port == x.Port && s.IP.Equal(x.IP)
}





// r.cli, err = etcdClientV3.New(etcdClientV3.Config{
	// 	Endpoints:   []string{"localhost:2379"},  
	// 	DialTimeout: 15 * time.Second,
	// 	LogConfig: &zap.Config{
	// 		Level:       zap.NewAtomicLevelAt(zap.ErrorLevel),
	// 		Development: false,
	// 		Sampling: &zap.SamplingConfig{
	// 			Initial:    100,
	// 			Thereafter: 100,
	// 		},
	// 		Encoding:      "json",
	// 		EncoderConfig: zap.NewProductionEncoderConfig(),
	// 		// Use "/dev/null" to discard all
	// 		OutputPaths:      []string{"stderr"},
	// 		ErrorOutputPaths: []string{"stderr"},
	// 	},
	// })