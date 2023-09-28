package terrors

var (
	// 0
	Succ = New(0, "success")

	// 600起始 default 默认
	ParamsEmpty = New(600000, "params empty") // 参数为空
	ParamsError = New(600001, "params error") // 参数错误
	FormatError = New(600002, "format error") // 参数格式错误
	AssertError = New(600003, "assert error") // 类型断言错误
	TransError  = New(600004, "trans error")  // 类型转换错误

	SystemError = New(600005, "system error") // 系统错误
	RPCError    = New(600006, "rpc error")    // 远程调用错误
	ConfError   = New(600007, "conf error")   // 内部配置错误
	NoData      = New(600008, "no data")      // 没有数据

	// 700起始 core 核心
	NameRegistered    = New(700000, "the name has registered") // 容器已注册
	NameNotRegistered = New(700001, "the name not regisered")  // 容器未注册

	// 800起始 user

	// 900起始 third

)
