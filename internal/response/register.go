package response

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type Register struct {
	DeviceId uint64
}

func (reg *Register) Decode(r io.Reader) (int, error) {
	dict := map[uint16]any{
		6: &reg.DeviceId,
	}
	return protocol.Decode(r, dict)
}

func (r Register) Type() uint16 {
	return 0x4000
}
