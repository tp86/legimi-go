package encoding

import "io"

type Array[T any] []T

func (a Array[T]) Encode(w io.Writer) error {
	err := Encode(w, uint16(len(a)))
	if err != nil {
		return err
	}
	for _, value := range a {
		err = Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a Array[T]) EncodedLength() int {
	length := U16Length
	for _, value := range a {
		length += EncodedLength(value)
	}
	return length
}

func (a *Array[T]) Decode(r io.Reader) (int, error) {
	var bytesRead int
	var count uint16
	bytesRead, err := Decode(r, &count)
	if err != nil {
		return bytesRead, err
	}
	array := *a
	// ensure array has at least count capacity
	array = extend(array, count)
	for i := uint16(0); i < count; i++ {
		var value T
		n, err := Decode(r, &value)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
		array[i] = value
	}
	*a = array
	return bytesRead, err
}

func extend[T any](a []T, count uint16) []T {
	diff := int(count) - cap(a)
	if diff > 0 {
		a = append(make([]T, count), a...)
	}
	return a
}
