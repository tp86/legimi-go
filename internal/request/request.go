package request

import "github.com/tp86/legimi-go/internal/protocol"

type Request interface {
	protocol.Encoder
	protocol.Typed
}
