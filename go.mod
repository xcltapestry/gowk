module github.com/xcltapestry/gowk

go 1.15

require (
	github.com/bsm/redislock v0.7.0
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4 // indirect
	github.com/go-redis/redis/v8 v8.4.10
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.5 // indirect
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/json-iterator/go v1.1.7
	github.com/manifoldco/promptui v0.8.0 // indirect
	github.com/prometheus/client_golang v1.0.0 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/viper v1.7.1
	go.etcd.io/bbolt v1.3.3 // indirect
	go.etcd.io/etcd v3.3.22+incompatible
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

// backend.go:324:9: options.OpenFile undefined (type bbolt.Options has no field or method OpenFile)
//replace github.com/coreos/bbolt => go.etcd.io/bbolt v3.3.20
// replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

// replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

// require (
// 	github.com/coreos/etcd v2.3.8+incompatible // indirect
// 	go.etcd.io/etcd v3.3.22+incompatible
// )

// replace google.golang.org/grpc => google.golang.org/grpc v1.26.0 // indirect

// 	github.com/coreos/etcd v3.3.22+incompatible
// v3.3.22+incompat
