package protocol

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
	case map[Key]any:
		return decodeMap(r, value)
	case *uint8:
		return U8Length, decode(r, value)
	case *uint16:
		return U16Length, decode(r, value)
	case *uint32:
		return U32Length, decode(r, value)
	case *uint64:
		return U64Length, decode(r, value)
	default:
		return 0, fmt.Errorf("decoding in unsupported for type %T", value)
	}
}

func DecodeWithLength(r io.Reader, value any, length int) (int, error) {
	switch value := value.(type) {
	case *string:
		bytes := make([]byte, length)
		bytesRead, err := r.Read(bytes)
		if err != nil {
			return bytesRead, err
		}
		*value = string(bytes)
		return bytesRead, nil
	default:
		return Decode(r, value)
	}
}

func SkipDecode(r io.Reader, n int) (int, error) {
	return r.Read(make([]byte, n))
}

func decode(r io.Reader, value any) error {
	return binary.Read(r, binary.LittleEndian, value)
}

func decodeMap(r io.Reader, value map[Key]any) (int, error) {
	var bytesRead int
	var count uint16
	bytesRead, err := Decode(r, &count)
	if err != nil {
		return bytesRead, err
	}
	for i := uint16(0); i < count; i++ {
		var key uint16
		n, err := Decode(r, &key)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
		if target, ok := value[key]; ok {
			n, err := Decode(r, WithLength{Value: target})
			bytesRead += n
			if err != nil {
				return bytesRead, err
			}
		} else {
			var toSkip uint32
			n, err := Decode(r, &toSkip)
			bytesRead += n
			if err != nil {
				return bytesRead, err
			}
			n, err = SkipDecode(r, int(toSkip))
			bytesRead += n
			if err != nil {
				return bytesRead, err
			}
		}
	}
	return bytesRead, err
}
