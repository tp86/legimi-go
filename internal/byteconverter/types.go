package byteconverter

import (
	"io"
	"reflect"
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
	val := uint64(n.val)
	for i := uintptr(0); i < n.size(); i, val = i+1, val>>8 {
		w.WriteByte(byte(val))
	}
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

type (
	bbyte  = numWithLength[uint8, length]
	bshort = numWithLength[uint16, length]
	bint   = numWithLength[uint32, length]
	blong  = numWithLength[uint64, length]
)

func (nl numWithLength[N, L]) writeBytesTo(w io.ByteWriter) {
	n := num[N]{nl.val}
	num[L]{L(n.size())}.writeBytesTo(w)
	n.writeBytesTo(w)
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

type sequence []ByteWriterTo

func (seq sequence) writeBytesTo(w io.ByteWriter) {
	for _, bwt := range seq {
		bwt.writeBytesTo(w)
	}
}
