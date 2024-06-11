package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/tp86/legimi-go/internal/api/protocol"
	"github.com/tp86/legimi-go/internal/api/protocol/encoding"
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
		3: protocol.AppVersion,
		4: uint32(0),
		5: uint64(0),
	}
}

type Session struct {
	state               uint16
	errorCode           uint16
	lastAgreementPeriod bool
	serverMessage       string
	maxTablets          uint32
	maxSmartphones      uint32
	maxEinks            uint32
	maxWins             uint32
	maxKindles          uint32
	maxDownloads        uint32
	downloadsLeft       uint32
	activating          bool
	maxDevices          uint32
	Id                  string
}

func (s *Session) Decode(r io.Reader) (int, error) {
	return encoding.Decode(r, encoding.Map{
		0:  &s.state,
		5:  &s.errorCode,
		37: &s.lastAgreementPeriod,
		3:  &s.serverMessage,
		15: &s.maxTablets,
		14: &s.maxSmartphones,
		16: &s.maxEinks,
		17: &s.maxWins,
		32: &s.maxKindles,
		33: &s.maxDownloads,
		34: &s.downloadsLeft,
		21: &s.activating,
		27: &s.maxDevices,
		7:  &s.Id,
	})
}

func (s Session) Type() uint16 {
	return 0x4002
}

func (s Session) DebugFormat() string {
	var b strings.Builder
	b.WriteString("Session response:\n")
	b.WriteString(fmt.Sprintf("state: %d\n", s.state))
	b.WriteString(fmt.Sprintf("error code: %d\n", s.errorCode))
	b.WriteString(fmt.Sprintf("in last agreement period: %t\n", s.lastAgreementPeriod))
	b.WriteString(fmt.Sprintf("is activating: %t\n", s.activating))
	b.WriteString(fmt.Sprintf(`server message: "%s"`+"\n", s.serverMessage))
	b.WriteString(fmt.Sprintf("max tablets: %d\n", s.maxTablets))
	b.WriteString(fmt.Sprintf("max smartphones: %d\n", s.maxSmartphones))
	b.WriteString(fmt.Sprintf("max einks: %d\n", s.maxEinks))
	b.WriteString(fmt.Sprintf("max windows: %d\n", s.maxWins))
	b.WriteString(fmt.Sprintf("max kindles: %d\n", s.maxKindles))
	b.WriteString(fmt.Sprintf("max devices: %d\n", s.maxDevices))
	b.WriteString(fmt.Sprintf("max downloads: %d\n", s.maxDownloads))
	b.WriteString(fmt.Sprintf("downloads left: %d", s.downloadsLeft))
	return b.String()
}
