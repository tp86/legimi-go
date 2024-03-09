package byteconverter_

import (
	"io"
	"reflect"
)

func newNumWithLength[N numbers](n N) value {
	return &numWithLength[N, length]{n}
}

// func newbstring[L lengths](s string) value {
//	return bstring[L]{s}
// }
//
// func newsequence(vals ...value) value {
//	return sequence(vals)
// }

var (
	Byte  = newNumWithLength[uint8]
	Short = newNumWithLength[uint16]
	Int   = newNumWithLength[uint32]
	Long  = newNumWithLength[uint64]
	// String      = newbstring[length]
	// ShortString = newbstring[shortLength]
	// Sequence    = newsequence
)

type numbers interface {
	uint8 | uint16 | uint32 | uint64
}

type num[N numbers] struct {
	val N
}

func (n num[N]) size() uintptr {
	return reflect.TypeFor[N]().Size()
}

func (n num[N]) writeBytesTo(w io.ByteWriter) {
	value := uint64(n.val)
	for i := uintptr(0); i < n.size(); i, value = i+1, value>>8 {
		w.WriteByte(byte(value))
	}
}

func (n *num[N]) readBytesFrom(r io.ByteReader) {
	value := uint64(0)
	for i := 0; i < 4; i, value = i+1, value<<8 {
		b, _ := r.ReadByte()
		value += uint64(b)
	}
	n.val = N(value)
}

type (
	length      = uint32
	shortLength = uint16
)
type lengths interface {
	length | shortLength
}

type numWithLength[N numbers, L lengths] struct {
	val N
}

func (nl numWithLength[N, L]) writeBytesTo(w io.ByteWriter) {
	n := num[N]{nl.val}
	num[L]{L(n.size())}.writeBytesTo(w)
	n.writeBytesTo(w)
}

func (nl *numWithLength[N, L]) readBytesFrom(r io.ByteReader) {
	l := num[L]{}
	l.readBytesFrom(r)
	nl.readBytesFrom(r)
}

type bstring[L lengths] struct {
	val string
}

func (s bstring[L]) writeBytesTo(w io.ByteWriter) {
	str := s.val
	strlen := len(str)
	l := num[L]{L(strlen)}
	l.writeBytesTo(w)
	for i := 0; i < strlen; i++ {
		w.WriteByte(str[i])
	}
}

type sequence []value

func (seq sequence) writeBytesTo(w io.ByteWriter) {
	for _, bwt := range seq {
		bwt.writeBytesTo(w)
	}
}

type value interface {
	ByteWriterTo
	ByteReaderFrom
}
