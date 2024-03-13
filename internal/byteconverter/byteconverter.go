package byteconverter

import (
	"errors"
	"io"
	"log"
	"reflect"
	"strings"
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

func readByte(r io.ByteReader) byte {
	b, err := r.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Fatalln("unexpected EOF while reading")
		}
	}
	return b
}

func (n *number[N]) readBytesFrom(r io.ByteReader) {
	var value uint64
	for i := uintptr(0); i < reflect.TypeOf(n.Value).Size(); i++ {
		value += uint64(readByte(r)) << (i * 8)
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
		readByte(r)
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

type bstring[L lengths] struct {
	Value string
}

func (bs *bstring[L]) writeBytesTo(w io.ByteWriter) {
	valueLength := len(bs.Value)
	l := number[L]{L(valueLength)}
	l.writeBytesTo(w)
	for i := 0; i < valueLength; i++ {
		w.WriteByte(bs.Value[i])
	}
}

func (bs *bstring[L]) readBytesFrom(r io.ByteReader) {
	l := &number[L]{}
	l.readBytesFrom(r)
	s := strings.Builder{}
	for i := 0; i < int(l.Value); i++ {
		s.WriteByte(readByte(r))
	}
	bs.Value = s.String()
}

type (
	String      = bstring[uint32]
	ShortString = bstring[uint16]
)

type bytesReaderWriter interface {
	bytesReaderFrom
	bytesWriterTo
}

type Sequence []bytesReaderWriter

func (s *Sequence) writeBytesTo(w io.ByteWriter) {
	for _, v := range *s {
		v.writeBytesTo(w)
	}
}

func (s *Sequence) readBytesFrom(r io.ByteReader) {
	for _, v := range *s {
		v.readBytesFrom(r)
	}
}

func WriteBytesTo(w io.ByteWriter, b bytesWriterTo) {
	b.writeBytesTo(w)
}

func ReadBytesFrom(r io.ByteReader, b bytesReaderFrom) {
	b.readBytesFrom(r)
}

type ToBytesMapper interface {
	MapToBytes() bytesWriterTo
}

type FromBytesMapper interface {
	MapFromBytes() bytesReaderFrom
}
