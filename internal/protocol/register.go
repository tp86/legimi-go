package protocol

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol/encoding"
)

type RegisterRequest struct {
	loginData
	kindleSerialNo string
}

func NewRegisterRequest(login, password, kindleSerialNo string) RegisterRequest {
	return RegisterRequest{
		loginData:      loginData{login: login, password: password},
		kindleSerialNo: fmt.Sprintf("Kindle||Kindle||%s||Kindle", kindleSerialNo),
	}
}

func (r RegisterRequest) Encode(w io.Writer) error {
	err := encoding.Encode(w, emptyId)
	if err != nil {
		return err
	}
	for _, value := range r.data() {
		err := encoding.Encode(w, uint16(stringLength(value)))
		if err != nil {
			return err
		}
		err = encoding.Encode(w, value)
		if err != nil {
			return err
		}
	}
	err = encoding.Encode(w, emptyLocale)
	if err != nil {
		return err
	}
	return nil
}

func (r RegisterRequest) EncodedLength() int {
	var totalLength int
	totalLength += encoding.EncodedLength(emptyId)
	for _, s := range r.data() {
		totalLength += encoding.U16Length + stringLength(s)
	}
	totalLength += encoding.EncodedLength(emptyLocale)
	return totalLength
}

func (r RegisterRequest) Type() uint16 {
	return 0x0042
}

func stringLength(s string) int {
	return encoding.EncodedLength(s)
}

var (
	emptyId     uint64 = 0
	emptyLocale        = encoding.Array[uint8]{}
)

func (r RegisterRequest) data() []string {
	return []string{r.login, r.password, r.kindleSerialNo}
}

type Register struct {
	KindleId uint64
}

func (reg *Register) Decode(r io.Reader) (int, error) {
	dict := encoding.Map{
		6: &reg.KindleId,
	}
	return encoding.Decode(r, dict)
}

func (r Register) Type() uint16 {
	return 0x4000
}
