package request

import (
	"bytes"
	"fmt"

	bc "github.com/tp86/legimi-go/internal/byteconverter"
	"github.com/tp86/legimi-go/internal/packet"
)

type Registration struct {
	Login          string
	Password       string
	KindleSerialNo string
}

func wrapSerialNo(serialNo string) string {
	return fmt.Sprintf("Kindle||Kindle||%s||Kindle", serialNo)
}

func (r Registration) Type() packet.Type {
	return RegistrationRequest
}

func (r Registration) ToBytes() []byte {
	buf := new(bytes.Buffer)
	(&bc.Sequence{
		&bc.RawLong{Value: 0},
		&bc.ShortString{Value: r.Login},
		&bc.ShortString{Value: r.Password},
		&bc.ShortString{Value: wrapSerialNo(r.KindleSerialNo)},
		&bc.RawShort{Value: 0},
	}).WriteBytesTo(buf)
	return buf.Bytes()
}
