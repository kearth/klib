package server

type Server interface {
}

func (s *Server) Run() error {
	return nil
}



interface Transport {
	// 初始化通道
	start(): Promise<void>;
	// 发送信息或响应
	send(message protocol.Notification) error
	
	close(): Promise<void>;
  
	onclose?: () => void;
  
	onerror?: (error: Error) => void;
	// 处理收到的请求或响应
	onmessage?: (message: JSONRPCMessage) => void;
	// 连接对应的回话 ID
	sessionId?: string;
  }