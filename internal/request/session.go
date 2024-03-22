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
	return protocol.Encode(w, s.asMap())
}

func (s Session) EncodedLength() int {
	return protocol.EncodedLength(s.asMap())
}

func (s Session) Type() uint16 {
	return 0x0050
}

func (s Session) asMap() protocol.Map {
	return protocol.Map{
		0: s.login,
		1: s.password,
		2: s.kindleId,
		3: protocol.AppVersion,
		4: uint32(0),
		5: uint64(0),
	}
}
