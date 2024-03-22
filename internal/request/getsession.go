package request

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type GetSession struct {
	loginData
	kindleId uint64
}

func NewGetSessionRequest(login, password string, kindleId uint64) GetSession {
	return GetSession{
		loginData: loginData{login: login, password: password},
		kindleId:  kindleId,
	}
}

func (s GetSession) Encode(w io.Writer) error {
	return protocol.Encode(w, s.asMap())
}

func (s GetSession) EncodedLength() int {
	return protocol.EncodedLength(s.asMap())
}

func (s GetSession) Type() uint16 {
	return 0x0050
}

func (s GetSession) asMap() protocol.Map {
	return protocol.Map{
		0: s.login,
		1: s.password,
		2: s.kindleId,
		3: protocol.AppVersion,
		4: uint32(0),
		5: uint64(0),
	}
}
