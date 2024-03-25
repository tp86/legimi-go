package protocol

import "github.com/tp86/legimi-go/internal/api/protocol/encoding"

type Request interface {
	encoding.Encoder
	encoding.Typed
}
