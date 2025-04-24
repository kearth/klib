package protocol

// Response 响应结构体
type Response struct {
	JSONRPC string         `json:"jsonrpc"`          // 必须为 "2.0"
	ID      any            `json:"id"`               // string
	Result  map[string]any `json:"result,omitempty"` // 调用成功时返回的结果
	Error   *RPCError      `json:"error,omitempty"`  // 调用失败时返回的错误信息
}

// RPCError 错误信息结构体
type RPCError struct {
	Code    int    `json:"code"`           // 错误码
	Message string `json:"message"`        // 错误信息
	Data    any    `json:"data,omitempty"` // 可选的附加数据
}
