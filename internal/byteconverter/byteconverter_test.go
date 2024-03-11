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
}

func TestWritingBytes(t *testing.T) {
	for tcname, tc := range writingTests {
		name := "writing " + tcname
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			WriteBytesTo(buf, tc.input)
			bs := buf.Bytes()
			if !slices.Equal(bs, tc.expected) {
				t.Errorf("%s: expectedn%v, got %v", name, tc.expected, bs)
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

var readingTests = map[string]struct {
	input    []byte
	output   value
	expected any
}{
	"raw byte":  {[]byte{1}, &RawByte{}, byte(1)},
	"raw short": {[]byte{6, 0}, &RawShort{}, uint16(6)},
	"raw int":   {[]byte{0x02, 0x40, 0, 0}, &RawInt{}, uint32(16386)},
	"raw long":  {[]byte{255, 255, 255, 255, 255, 255, 255, 0x7f}, &RawLong{}, uint64(9223372036854775807)},
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
