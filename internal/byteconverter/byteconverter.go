package byteconverter

import (
	"io"
	"reflect"
)

type numbers interface {
	uint8 | uint16 | uint32 | uint64
}

type bytesReaderFrom interface {
	readBytesFrom(r io.ByteReader)
}

type bytesWriterTo interface {
	writeBytesTo(w io.ByteWriter)
}

type number[N numbers] struct {
	Value N
}

func (n *number[N]) writeBytesTo(w io.ByteWriter) {
	value := uint64(n.Value)
	for i := uintptr(0); i < reflect.TypeOf(n.Value).Size(); i, value = i+1, value>>8 {
		w.WriteByte(byte(value))
	}
}

func (n *number[N]) readBytesFrom(r io.ByteReader) {
	var value uint64
	for i := uintptr(0); i < reflect.TypeOf(n.Value).Size(); i++ {
		b, _ := r.ReadByte()
		value += uint64(b) << (i * 8)
	}
	n.Value = N(value)
}

type (
	RawByte  = number[uint8]
	RawShort = number[uint16]
	RawInt   = number[uint32]
	RawLong  = number[uint64]
)

type lengths interface {
	uint16 | uint32
}

type numberWithLength[N numbers, L lengths] struct {
	Value N
}

func (nl *numberWithLength[N, L]) writeBytesTo(w io.ByteWriter) {
	l := number[L]{L(reflect.TypeOf(nl.Value).Size())}
	l.writeBytesTo(w)
	n := number[N]{nl.Value}
	n.writeBytesTo(w)
}

func (nl *numberWithLength[N, L]) readBytesFrom(r io.ByteReader) {
	for i := uintptr(0); i < reflect.TypeOf(L(0)).Size(); i++ {
		r.ReadByte()
	}
	n := number[N]{}
	n.readBytesFrom(r)
	nl.Value = n.Value
}

type (
	Byte  = numberWithLength[uint8, uint32]
	Short = numberWithLength[uint16, uint32]
	Int   = numberWithLength[uint32, uint32]
	Long  = numberWithLength[uint64, uint32]
)

func WriteBytesTo(w io.ByteWriter, b bytesWriterTo) {
	b.writeBytesTo(w)
}

func ReadBytesFrom(r io.ByteReader, b bytesReaderFrom) {
	b.readBytesFrom(r)
}
