package protocol

import "io"

const (
	Version    = uint32(17)
	AppVersion = "1.8.5 Windows"
)

const (
	U8Length  = 1
	U16Length = 2
	U32Length = 4
	U64Length = 8
)

type Key = uint16

type Typed interface {
	Type() uint16
}

type WithLength struct {
	Value any
}

func (wl WithLength) Encode(w io.Writer) error {
	err := Encode(w, uint32(EncodedLength(wl.Value)))
	if err != nil {
		return err
	}
	return Encode(w, wl.Value)
}

func (wl WithLength) EncodedLength() int {
	return U32Length + EncodedLength(wl.Value)
}

func (wl WithLength) Decode(r io.Reader) (int, error) {
	var length uint32
	bytesRead, err := Decode(r, &length)
	if err != nil {
		return bytesRead, err
	}
	n, err := DecodeWithLength(r, wl.Value, int(length))
	bytesRead += n
	return bytesRead, err
}
