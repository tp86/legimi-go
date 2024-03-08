package byteconverter

import (
	"io"
	"reflect"

	"golang.org/x/exp/constraints"
)

type ByteWriterTo interface {
	writeBytesTo(w io.ByteWriter)
}

func WriteBytesTo(w io.ByteWriter, val ByteWriterTo) {
	val.writeBytesTo(w)
}

type num[I constraints.Integer] struct {
	val I
}

type (
	shortLength = num[uint16]
	length      = num[uint32]
)

type numWithLength[I constraints.Integer] struct {
	val I
}

type (
	bbyte  = numWithLength[uint8]
	bshort = numWithLength[int16]
	bint   = numWithLength[int32]
	blong  = numWithLength[int64]
)

func (n num[I]) writeBytesTo(w io.ByteWriter) {
	val := uint64(n.val)
	for i := uintptr(0); i < reflect.TypeOf(n.val).Size(); i, val = i+1, val>>8 {
		w.WriteByte(byte(val))
	}
}

func (nl numWithLength[I]) writeBytesTo(w io.ByteWriter) {
	length{uint32(reflect.TypeOf(nl.val).Size())}.writeBytesTo(w)
	num[I]{nl.val}.writeBytesTo(w)
}

func Byte(b uint8) ByteWriterTo {
	return bbyte{val: b}
}

func Short(s int16) ByteWriterTo {
	return bshort{val: s}
}

func Int(i int32) ByteWriterTo {
	return bint{val: i}
}

func Long(l int64) ByteWriterTo {
	return blong{val: l}
}

type bstring[L num[LI], LI uint16 | uint32] struct {
	val string
}

func asLen[L num[LI], LI uint16 | uint32](i int) num[LI] {
	return num[LI]{val: LI(i)}
}

func (s bstring[L, LI]) writeBytesTo(w io.ByteWriter) {
	str := s.val
	l := asLen[L, LI](len(str))
	l.writeBytesTo(w)
	for i := 0; i < len(str); i++ {
		w.WriteByte(str[i])
	}
}

func String(s string) ByteWriterTo {
	return bstring[num[uint32], uint32]{val: s}
}

func ShortString(s string) ByteWriterTo {
	return bstring[num[uint16], uint16]{val: s}
}

type sequence []ByteWriterTo

func (seq sequence) writeBytesTo(w io.ByteWriter) {
	for _, bwt := range seq {
		bwt.writeBytesTo(w)
	}
}

func Sequence(bvals ...ByteWriterTo) ByteWriterTo {
	return sequence(bvals)
}
