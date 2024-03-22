package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Encoder interface {
	Encode(w io.Writer) error
	EncodedLength() int
}

func Encode(w io.Writer, value any) error {
	switch value := value.(type) {
	case Encoder:
		return value.Encode(w)
	case uint8, uint16, uint32, uint64:
		return encode(w, value)
	case bool:
		var byteValue uint8
		if value {
			byteValue = 0xff
		}
		return encode(w, byteValue)
	case string:
		return encode(w, []byte(value))
	default:
		return fmt.Errorf("encoding is unsupported for type %T", value)
	}
}

func EncodedLength(value any) int {
	switch value := value.(type) {
	case Encoder:
		return value.EncodedLength()
	case uint8, bool, *uint8, *bool:
		return U8Length
	case uint16, *uint16:
		return U16Length
	case uint32, *uint32:
		return U32Length
	case uint64, *uint64:
		return U64Length
	case string:
		return len(value)
	default:
		return 0
	}
}

func encode(w io.Writer, value any) error {
	return binary.Write(w, binary.LittleEndian, value)
}
