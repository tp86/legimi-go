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
	case []uint8:
		return encodeArray(w, value)
	case []uint16:
		return encodeArray(w, value)
	case []uint32:
		return encodeArray(w, value)
	case []uint64:
		return encodeArray(w, value)
	case []Encoder:
		return encodeArray(w, value)
	case []string:
		return encodeArray(w, value)
	case []any:
		return encodeArray(w, value)
	case uint8, uint16, uint32, uint64:
		return encode(w, value)
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
	case []Encoder:
		return arrayLength(value)
	case []uint8:
		return arrayLength(value)
	case []uint16:
		return arrayLength(value)
	case []uint32:
		return arrayLength(value)
	case []uint64:
		return arrayLength(value)
	case []string:
		return arrayLength(value)
	case []any:
		return arrayLength(value)
	case uint8:
		return U8Length
	case uint16:
		return U16Length
	case uint32:
		return U32Length
	case uint64:
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

func encodeArray[T any](w io.Writer, values []T) error {
	err := encode(w, uint16(len(values)))
	if err != nil {
		return err
	}
	for _, value := range values {
		err = Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func arrayLength[T any](values []T) int {
	length := U16Length
	for _, value := range values {
		length += EncodedLength(value)
	}
	return length
}
