package protocol

// Request 请求结构体
type Request struct {
	JSONRPC string         `json:"jsonrpc"`          // 必须为 "2.0"
	ID      string         `json:"id"`               // 必须为字符串
	Method  string         `json:"method"`           // 方法名
	Params  map[string]any `json:"params,omitempty"` // 可选参数，key-value 格式
}
