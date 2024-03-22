package request

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

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
	err := protocol.Encode(w, emptyId)
	if err != nil {
		return err
	}
	for _, value := range r.data() {
		err := protocol.Encode(w, uint16(stringLength(value)))
		if err != nil {
			return err
		}
		err = protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	err = protocol.Encode(w, emptyLocale)
	if err != nil {
		return err
	}
	return nil
}

func (r Register) EncodedLength() int {
	var totalLength int
	totalLength += protocol.EncodedLength(emptyId)
	for _, s := range r.data() {
		totalLength += protocol.U16Length + stringLength(s)
	}
	totalLength += protocol.EncodedLength(emptyLocale)
	return totalLength
}

func (r Register) Type() uint16 {
	return 0x0042
}

func stringLength(s string) int {
	return protocol.EncodedLength(s)
}

var (
	emptyId     uint64 = 0
	emptyLocale        = protocol.Array[uint8]{}
)

func (r Register) data() []string {
	return []string{r.login, r.password, r.kindleSerialNo}
}
