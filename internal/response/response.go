package response

import "github.com/tp86/legimi-go/internal/protocol"

type Response interface {
	protocol.Decoder
	protocol.Typed
}
