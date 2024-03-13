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

type BytesReaderFrom interface {
	ReadBytesFrom(r io.ByteReader)
}

type BytesWriterTo interface {
	WriteBytesTo(w io.ByteWriter)
}

type number[N numbers] struct {
	Value N
}

func (n *number[N]) WriteBytesTo(w io.ByteWriter) {
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

func (n *number[N]) ReadBytesFrom(r io.ByteReader) {
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

func (nl *numberWithLength[N, L]) WriteBytesTo(w io.ByteWriter) {
	l := number[L]{L(reflect.TypeOf(nl.Value).Size())}
	l.WriteBytesTo(w)
	n := number[N]{nl.Value}
	n.WriteBytesTo(w)
}

func skipBytes(r io.ByteReader, bytesToSkip int) {
	for i := 0; i < bytesToSkip; i++ {
		readByte(r)
	}
}

func (nl *numberWithLength[N, L]) ReadBytesFrom(r io.ByteReader) {
	skipBytes(r, int(reflect.TypeOf(L(0)).Size()))
	n := number[N]{}
	n.ReadBytesFrom(r)
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

func (bs *bstring[L]) WriteBytesTo(w io.ByteWriter) {
	valueLength := len(bs.Value)
	l := number[L]{L(valueLength)}
	l.WriteBytesTo(w)
	for i := 0; i < valueLength; i++ {
		w.WriteByte(bs.Value[i])
	}
}

func (bs *bstring[L]) ReadBytesFrom(r io.ByteReader) {
	l := &number[L]{}
	l.ReadBytesFrom(r)
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
	BytesReaderFrom
	BytesWriterTo
}

type Bytes struct {
	Value []byte
}

func (rb *Bytes) WriteBytesTo(w io.ByteWriter) {
	(&RawInt{uint32(len(rb.Value))}).WriteBytesTo(w)
	for _, b := range rb.Value {
		w.WriteByte(b)
	}
}

func (rb *Bytes) ReadBytesFrom(r io.ByteReader) {
	l := &RawInt{}
	l.ReadBytesFrom(r)
	if rb.Value == nil {
		rb.Value = make([]byte, l.Value)
	}
	for i := uint32(0); i < l.Value; i++ {
		rb.Value[i] = readByte(r)
	}
}

type Dictionary struct {
	Value map[uint16]bytesReaderWriter
}

func (d *Dictionary) ReadBytesFrom(r io.ByteReader) {
	key := &RawShort{}
	key.ReadBytesFrom(r)
	count := key.Value
	length := &RawInt{}
	for i := uint16(0); i < count; i++ {
		key.ReadBytesFrom(r)
		if targetValue, ok := d.Value[key.Value]; ok {
			targetValue.ReadBytesFrom(r)
		} else {
			length.ReadBytesFrom(r)
			skipBytes(r, int(length.Value))
		}
	}
}

func WriteBytesTo(w io.ByteWriter, b BytesWriterTo) {
	b.WriteBytesTo(w)
}

func ReadBytesFrom(r io.ByteReader, b BytesReaderFrom) {
	b.ReadBytesFrom(r)
}

type bytesWriterSetter[V any] interface {
	newWriterTo(v V) BytesWriterTo
}

func (n number[N]) newWriterTo(v N) BytesWriterTo {
	return &number[N]{v}
}

func (nl numberWithLength[N, L]) newWriterTo(v N) BytesWriterTo {
	return &numberWithLength[N, L]{v}
}

func (bs bstring[L]) newWriterTo(v string) BytesWriterTo {
	return &bstring[L]{v}
}

func WriteAsBytesTo[V bytesWriterSetter[T], T any](w io.ByteWriter, v T) {
	var s V
	s.newWriterTo(v).WriteBytesTo(w)
}
