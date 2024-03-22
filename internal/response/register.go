package response

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type Register struct {
	KindleId uint64
}

func (reg *Register) Decode(r io.Reader) (int, error) {
	dict := protocol.Map{
		6: &reg.KindleId,
	}
	return protocol.Decode(r, dict)
}

func (r Register) Type() uint16 {
	return 0x4000
}
