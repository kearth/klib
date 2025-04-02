package resources

// Request 定义请求结构体
type Request struct {
	Method string `json:"method"`
}

// Response 定义响应结构体
type Response struct {
	URI         string  `json:"uri"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	MimeType    *string `json:"mimeType,omitempty"`
}
