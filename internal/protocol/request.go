package protocol

import "github.com/tp86/legimi-go/internal/protocol/encoding"

type Request interface {
	encoding.Encoder
	encoding.Typed
}
