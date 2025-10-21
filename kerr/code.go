package kerr

// 常用通用错误码常量
var (
	// =====================
	// 0 - 成功
	// =====================
	Succ = New(0, "success").WithDisplay("操作成功")

	// =====================
	// 1xxxx - 系统/基础设施错误
	// =====================
	SystemError     = New(10000, "system error").WithDisplay("系统错误，请稍后重试")
	ConfigError     = New(10001, "configuration error").WithDisplay("系统配置错误")
	TimeoutError    = New(10002, "operation timeout").WithDisplay("请求超时")
	PanicError      = New(10003, "panic occurred").WithDisplay("系统异常")
	InternalIOError = New(10004, "io error").WithDisplay("IO 操作失败")

	// =====================
	// 2xxxx - 用户/认证/权限错误
	// =====================
	Unauthorized  = New(20000, "unauthorized").WithDisplay("未授权，请登录")
	Forbidden     = New(20001, "forbidden").WithDisplay("没有访问权限")
	UserNotFound  = New(20002, "user not found").WithDisplay("用户不存在")
	PasswordError = New(20003, "invalid password").WithDisplay("密码错误")
	TokenExpired  = New(20004, "token expired").WithDisplay("登录状态已过期")
	AccountLocked = New(20005, "account locked").WithDisplay("账号已锁定")
	QuotaExceeded = New(20006, "quota exceeded").WithDisplay("已超出使用限制")

	// =====================
	// 3xxxx - 业务逻辑错误
	// =====================
	InvalidState      = New(30000, "invalid state").WithDisplay("当前状态不允许此操作")
	OperationConflict = New(30001, "operation conflict").WithDisplay("操作冲突或重复执行")
	DependencyMissing = New(30002, "dependency missing").WithDisplay("依赖缺失或未初始化")
	NoData            = New(30003, "no data").WithDisplay("没有符合条件的数据")
	ValidationFailed  = New(30004, "validation failed").WithDisplay("数据验证失败")

	// =====================
	// 4xxxx - 外部依赖错误
	// =====================
	DBError           = New(40000, "database error").WithDisplay("数据库异常")
	CacheError        = New(40001, "cache error").WithDisplay("缓存服务异常")
	NetworkError      = New(40002, "network error").WithDisplay("网络异常，请稍后再试")
	ThirdPartyError   = New(40003, "third-party service error").WithDisplay("第三方服务异常")
	MessageQueueError = New(40004, "message queue error").WithDisplay("消息队列异常")

	// =====================
	// 5xxxx - 核心组件/框架错误
	// =====================
	NameRegistered    = New(50000, "name already registered").WithDisplay("名称已注册")
	NameNotRegistered = New(50001, "name not registered").WithDisplay("名称未注册")
	ResourcePoolEmpty = New(50002, "resource pool empty").WithDisplay("资源池为空")
	ServiceNotReady   = New(50003, "service not ready").WithDisplay("服务未就绪")

	// =====================
	// 6xxxx - 预留扩展错误码
	// =====================
)
