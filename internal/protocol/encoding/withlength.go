package encoding

import "io"

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
	n, err := decodeWithLength(r, wl.Value, int(length))
	bytesRead += n
	return bytesRead, err
}

func decodeWithLength(r io.Reader, value any, length int) (int, error) {
	switch value := value.(type) {
	case *string:
		// special case for strings - we need length value to read correct number of bytes
		bytes := make([]byte, length)
		bytesRead, err := r.Read(bytes)
		if err != nil {
			return bytesRead, err
		}
		// just cast to string, no decoding needed
		*value = string(bytes)
		return bytesRead, nil
	default:
		return Decode(r, value)
	}
}
