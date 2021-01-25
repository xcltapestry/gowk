package main

import (
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	fmt.Println(" ---- main -----")

	cli := NewEtcdCli(
		WithAddress([]string{"WithPrefix"}),
		WithDialTimeout(2*time.Second),
		WithRequestTimeout(1*time.Second),
		WithPrefix("/test/"))

	err := cli.Connect()
	if err != nil {
		fmt.Println(" err:", err)
		return
	}
	defer cli.Close()

	cli.ReadfileToETCD()
}

type EtcdCli struct {
	config         clientv3.Config
	prefix         string
	requestTimeout time.Duration

	client *clientv3.Client
}

type EtcdCliOption func(*EtcdCli)

func NewEtcdCli(options ...func(*EtcdCli)) *EtcdCli {

	cli := &EtcdCli{}
	cli.config.Endpoints = []string{"localhost:2379"}
	cli.config.DialTimeout = 2 * time.Second

	for _, f := range options {
		f(cli)
	}

	return cli

}

func WithAddress(addrs []string) EtcdCliOption {
	return func(c *EtcdCli) {
		c.config.Endpoints = addrs
		fmt.Println("[WithAddress] addrs:", addrs)
	}
}

func WithDialTimeout(dialTimeout time.Duration) EtcdCliOption {
	return func(c *EtcdCli) {

		c.config.DialTimeout = dialTimeout
		fmt.Println("[WithDialTimeout] dialTimeout:", dialTimeout)
	}
}

func WithPrefix(prefix string) EtcdCliOption {
	return func(c *EtcdCli) {
		c.prefix = prefix
		fmt.Println("[WithPrefix] prefix:", prefix)
	}
}

func WithRequestTimeout(requestTimeout time.Duration) EtcdCliOption {
	return func(c *EtcdCli) {
		c.requestTimeout = requestTimeout
		fmt.Println("[WithRequestTimeout] requestTimeout:", requestTimeout)
	}
}

func (e *EtcdCli) Connect() error {
	fmt.Println("[Connect] ---")

	var err error
	e.client, err = clientv3.New(e.config)
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdCli) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *EtcdCli) ReadfileToETCD() {
	fmt.Println("[ReadfileToETCD] ---")
}

func (e *EtcdCli) LoadfileFromETCD() {
	fmt.Println("[LoadfileFromETCD] ---")
}

func (e *EtcdCli) BuildConfig() {
	fmt.Println("[BuildConfig] ---")
}

/*

go: finding module for package github.com/coreos/go-systemd/journal
../../../../../../worklibs/gopath/pkg/mod/github.com/coreos/etcd@v3.3.25+incompatible/pkg/logutil/zap_journal.go:29:2: no matching versions for query "latest"

https://github.com/etcd-io/etcd/issues/11345

$ mkdir -p $GOPATH/src/github.com/coreos/go-systemd/
$ git clone git@github.com:coreos/go-systemd.git $GOPATH/src/github.com/coreos/go-systemd/
$ cd $myproject
$ go mod edit -replace github.com/coreos/go-systemd/=/Users/dselans/Code/go/src/github.com/coreos/go-systemd
$ go get go.etcd.io/etcd/clientv3

# undefined: resolver.BuildOption

*/
