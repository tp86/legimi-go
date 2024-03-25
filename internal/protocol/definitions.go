package protocol

import "github.com/tp86/legimi-go/internal/protocol/encoding"

const (
	Version    = uint32(17)
	AppVersion = "1.8.5 Windows"
)

var emptyLength = encoding.WithLength{}.EncodedLength()
