package protocol

import "github.com/tp86/legimi-go/internal/protocol/encoding"

type Response interface {
	encoding.Decoder
	encoding.Typed
}
