package request

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type Request interface {
	protocol.Encoder
	Type() uint16
}

type loginData struct {
	login    string
	password string
}

type Register struct {
	loginData
	kindleSerialNo string
}

func NewRegisterRequest(login, password, kindleSerialNo string) Register {
	return Register{
		loginData:      loginData{login: login, password: password},
		kindleSerialNo: fmt.Sprintf("Kindle||Kindle||%s||Kindle", kindleSerialNo),
	}
}

func (r Register) Encode(w io.Writer) error {
	for _, value := range []any{
		uint64(0),
		uint16(len(r.login)),
		r.login,
		uint16(len(r.password)),
		r.password,
		uint16(len(r.kindleSerialNo)),
		r.kindleSerialNo,
		[]uint8{},
	} {
		err := protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r Register) EncodedLength() int {
	return protocol.U64Length +
		protocol.U16Length + protocol.EncodedLength(r.login) +
		protocol.U16Length + protocol.EncodedLength(r.password) +
		protocol.U16Length + protocol.EncodedLength(r.kindleSerialNo) +
		protocol.EncodedLength([]uint8{})
}

func (r Register) Type() uint16 {
	return 0x0042
}
