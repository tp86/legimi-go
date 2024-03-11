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
	"raw byte": {&RawByte{1}, []byte{1}},
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
	"raw byte": {[]byte{1}, &RawByte{}, byte(1)},
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
