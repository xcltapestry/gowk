package naming

var _eRegister  *etcdRegister

func init(){
	_eRegister = NewEtcdRegister()
}

func Register(serviceName,addr string)error {
	 return _eRegister.Register(serviceName,addr)
}

func CloseRegisterEtcd(){
	_eRegister.CloseRegisterEtcd()
}




