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
	bc.WriteAsBytesTo[bc.RawLong](buf, 0)
	bc.WriteAsBytesTo[bc.ShortString](buf, r.Login)
	bc.WriteAsBytesTo[bc.ShortString](buf, r.Password)
	bc.WriteAsBytesTo[bc.ShortString](buf, wrapSerialNo(r.KindleSerialNo))
	bc.WriteAsBytesTo[bc.RawShort](buf, 0)
	return buf.Bytes()
}
