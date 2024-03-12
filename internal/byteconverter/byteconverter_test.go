package byteconverter

import (
	"bytes"
	"slices"
	"testing"
)

var writingTests = map[string]struct {
	input    bytesWriterTo
	expected []byte
}{
	"raw byte":  {&RawByte{1}, []byte{1}},
	"raw short": {&RawShort{6}, []byte{6, 0}},
	"raw int":   {&RawInt{16386}, []byte{0x02, 0x40, 0, 0}},
	"raw long":  {&RawLong{9223372036854775807}, []byte{255, 255, 255, 255, 255, 255, 255, 0x7f}},
	"byte":      {&Byte{17}, []byte{1, 0, 0, 0, 0x11}},
	"short":     {&Short{4}, []byte{2, 0, 0, 0, 4, 0}},
	"int":       {&Int{24}, []byte{4, 0, 0, 0, 0x18, 0, 0, 0}},
	"long":      {&Long{8}, []byte{8, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0}},
}

func TestWritingBytes(t *testing.T) {
	for tcname, tc := range writingTests {
		name := "writing " + tcname
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			WriteBytesTo(buf, tc.input)
			bs := buf.Bytes()
			if !slices.Equal(bs, tc.expected) {
				t.Errorf("%s: expected %v, got %v", name, tc.expected, bs)
			}
		})
	}
}

type value interface {
	bytesReaderFrom
	value() any
}

func (n number[N]) value() any {
	return n.Value
}

func (nl numberWithLength[N, L]) value() any {
	return nl.Value
}

var readingTests = map[string]struct {
	input    []byte
	output   value
	expected any
}{
	"raw byte":  {[]byte{1}, &RawByte{}, byte(1)},
	"raw short": {[]byte{6, 0}, &RawShort{}, uint16(6)},
	"raw int":   {[]byte{0x02, 0x40, 0, 0}, &RawInt{}, uint32(16386)},
	"raw long":  {[]byte{255, 255, 255, 255, 255, 255, 255, 0x7f}, &RawLong{}, uint64(9223372036854775807)},
	"byte":      {[]byte{1, 0, 0, 0, 1}, &Byte{}, byte(1)},
	"short":     {[]byte{2, 0, 0, 0, 3, 0}, &Short{}, uint16(3)},
	"int":       {[]byte{4, 0, 0, 0, 0x02, 0x40, 0, 0}, &Int{}, uint32(16386)},
	"long":      {[]byte{8, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 0x7f}, &Long{}, uint64(9223372036854775807)},
}

func TestReadingBytes(t *testing.T) {
	for tcname, tc := range readingTests {
		name := "reading " + tcname
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(tc.input)
			ReadBytesFrom(buf, tc.output)
			if tc.output.value() != tc.expected {
				t.Errorf("%s: expected: %v, got %v", name, tc.expected, tc.output.value())
			}
		})
	}
}
