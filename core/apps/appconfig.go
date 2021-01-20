package apps

func InitConfig(path, fileName string) ApplicationOption {
	return func(a *Application) {
		a.Config.ReadInConfig(path, fileName)
		a.buildConfig()
	}
}

// func (app *Application) SyncToETCD(yamlName string) error {

// 	endpoints := []string{"localhost:2379"}
// 	dialTimeout := 5 * time.Second

// 	etcdClient, err := confd.New(clientv3.Config{
// 		Endpoints:   endpoints,
// 		DialTimeout: dialTimeout,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	defer etcdClient.Close()
// 	fmt.Println(" etcd 连接成功。 ")

// 	for i, k := range app.Config.AllKeys() {
// 		fmt.Println(i, "==>", k)

// 	}

// 	// key1, value1 := "testkey1", "value"

// 	// ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
// 	// _, err = cli.Put(ctx, key1, value1)
// 	// cancel()
// 	// if err != nil {
// 	//     log.Println("Put failed. ", err)
// 	// } else {
// 	//     log.Printf("Put {%s:%s} succeed\n", key1, value1)
// 	// }

// 	return nil
// }
