package request

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type Session struct {
	loginData
	kindleId uint64
}

func NewSessionRequest(login, password string, kindleId uint64) Session {
	return Session{
		loginData: loginData{login: login, password: password},
		kindleId:  kindleId,
	}
}

func (s Session) Encode(w io.Writer) error {
	dict := map[protocol.Key]any{
		0: s.login,
		1: s.password,
		2: s.kindleId,
		3: protocol.AppVersion,
		4: uint32(0),
		5: uint64(0),
	}
	return protocol.Encode(w, dict)
}

func (s Session) EncodedLength() int {
	values := []int{
		protocol.EncodedLength(s.login),
		protocol.EncodedLength(s.password),
		protocol.EncodedLength(s.kindleId),
		protocol.EncodedLength(protocol.AppVersion),
		protocol.U32Length,
		protocol.U64Length,
	}
	return protocol.MapEncodedLength(values)
}

func (s Session) Type() uint16 {
	return 0x0050
}
