package encoding

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(r io.Reader) (int, error)
}

func Decode(r io.Reader, value any) (int, error) {
	switch value := value.(type) {
	case Decoder:
		return value.Decode(r)
	case *bool, *uint8, *uint16, *uint32, *uint64:
		return EncodedLength(value), decode(r, value)
	// case *string:
	//	return WithLength{Value: value}.Decode(r)
	default:
		return 0, fmt.Errorf("decoding is unsupported for type %T", value)
	}
}

func SkipDecode(r io.Reader, n int) (int, error) {
	return r.Read(make([]byte, n))
}

func SkipDecodeMany(r io.Reader, toSkip []int) (bytesSkip int, err error) {
	for _, skip := range toSkip {
		n, err := SkipDecode(r, skip)
		bytesSkip += n
		if err != nil {
			return bytesSkip, err
		}
	}
	return
}

func decode(r io.Reader, value any) error {
	return binary.Read(r, binary.LittleEndian, value)
}
