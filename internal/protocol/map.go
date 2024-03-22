package protocol

import "io"

type Map map[Key]any

func (m Map) Encode(w io.Writer) error {
	err := Encode(w, uint16(len(m)))
	if err != nil {
		return err
	}
	for key, value := range m {
		err = encode(w, key)
		if err != nil {
			return err
		}
		err = Encode(w, WithLength{Value: value})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Map) EncodedLength() int {
	totalLength := U16Length
	for _, value := range m {
		totalLength += U16Length /*key*/ + U32Length /*length*/ + EncodedLength(value)
	}
	return totalLength
}

func (m Map) Decode(r io.Reader) (int, error) {
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
		if target, ok := m[key]; ok {
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
