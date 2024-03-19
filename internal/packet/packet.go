package packet

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
	"github.com/tp86/legimi-go/internal/request"
)

type Packet struct {
	Req request.Request
}

func (p Packet) Encode(w io.Writer) error {
	for _, value := range []any{
		uint32(17),
		p.Req.Type(),
		protocol.WithLength{Value: p.Req},
	} {
		err := protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}
