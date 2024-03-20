package response

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type Session struct {
	Id string
}

func (s *Session) Decode(r io.Reader) (int, error) {
	dict := map[protocol.Key]any{
		7: &s.Id,
	}
	return protocol.Decode(r, dict)
}

func (s Session) Type() uint16 {
	return 0x4002
}
