package transport

import "github.com/kearth/klib/kmcp/protocol"

type Transport interface {
	Start() error
	Send(message protocol.Notification) error
	Close() error
	OnClose() error
	OnError(err error) error
	OnMessage(message protocol.Notification) error
	SessionID() string
}
