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
	err := encoding.Encode(w, emptyKindleId)
	if err != nil {
		return err
	}
	for _, value := range r.data() {
		err := encoding.Encode(w, uint16(encoding.EncodedLength(value)))
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
	totalLength += encoding.EncodedLength(emptyKindleId)
	for _, s := range r.data() {
		totalLength += encoding.U16Length + encoding.EncodedLength(s)
	}
	totalLength += encoding.EncodedLength(emptyLocale)
	return totalLength
}

func (r RegisterRequest) Type() uint16 {
	return 0x0042
}

var (
	emptyKindleId uint64 = 0
	emptyLocale          = encoding.Array[uint8]{}
)

func (r RegisterRequest) data() []string {
	return []string{r.login, r.password, r.kindleSerialNo}
}

type Register struct {
	KindleId uint64
}

func (reg *Register) Decode(r io.Reader) (int, error) {
	return encoding.Decode(r, encoding.Map{6: &reg.KindleId})
}

func (r Register) Type() uint16 {
	return 0x4000
}
