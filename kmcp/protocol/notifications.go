package protocol

// Notification 通知结构体
type Notification struct {
	JSONRPC string         `json:"jsonrpc"`          // 必须为 "2.0"
	Method  string         `json:"method"`           // 方法名
	Params  map[string]any `json:"params,omitempty"` // 参数（可选）
}
