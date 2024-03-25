package protocol

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol/encoding"
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
	return encoding.Encode(w, s.asMap())
}

func (s GetSession) EncodedLength() int {
	return encoding.EncodedLength(s.asMap())
}

func (s GetSession) Type() uint16 {
	return 0x0050
}

func (s GetSession) asMap() encoding.Map {
	return encoding.Map{
		0: s.login,
		1: s.password,
		2: s.kindleId,
		3: AppVersion,
		4: uint32(0),
		5: uint64(0),
	}
}

type Session struct {
	Id string
}

func (s *Session) Decode(r io.Reader) (int, error) {
	return encoding.Decode(r, encoding.Map{7: &s.Id})
}

func (s Session) Type() uint16 {
	return 0x4002
}
