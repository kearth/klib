package client

import (
	"github.com/kearth/klib/kmcp/protocol"
	"github.com/kearth/klib/kmcp/utils"
)

const (
	InitializeMethod = "initialize"
)

type InitializeRequest struct {
	protocol.Request
}

func Initialize() {
}

func NewInitializeRequest() *InitializeRequest {
	ir := &InitializeRequest{}
	ir.JSONRPC = protocol.JSONRPCVersion
	ir.Method = InitializeMethod
	ir.ID = utils.GenID()
	ir.Params = make(map[string]interface{})
	ir.Params["protocolVersion"] = protocol.ProtocalVersion
	ir.Params["capabilities"] = map[string]interface{}{
		"roots": map[string]interface{}{
			"listChanged": true,
		},
		"sampling": map[string]interface{}{},
	}
	ir.Params["clientInfo"] = map[string]interface{}{
		"name":    "ExampleClient",
		"version": "1.0.0",
	}
	return ir
}
