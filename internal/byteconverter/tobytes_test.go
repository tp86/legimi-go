package byteconverter_test

import (
	"bytes"
	"slices"
	"testing"

	bc "github.com/tp86/legimi-go/internal/byteconverter"
)

var convertToBytesTests = map[string]struct {
	input    bc.ByteWriterTo
	expected []byte
}{
	"byte":    {bc.Byte(254), []byte{1, 0, 0, 0, 0xfe}},
	"short":   {bc.Short(6), []byte{2, 0, 0, 0, 0x06, 0}},
	"integer": {bc.Int(17), []byte{4, 0, 0, 0, 0x11, 0, 0, 0}},
	"long":    {bc.Long(117928), []byte{8, 0, 0, 0, 0xa8, 0xcc, 0x01, 0, 0, 0, 0, 0}},
	"string":  {bc.String("abcdef"), []byte{6, 0, 0, 0, 'a', 'b', 'c', 'd', 'e', 'f'}},
	"string with short length": {
		input:    bc.ShortString("abcdef"),
		expected: []byte{6, 0, 'a', 'b', 'c', 'd', 'e', 'f'},
	},
	"sequence": {
		input:    bc.Sequence(bc.Int(0), bc.String("abc"), bc.Short(1024)),
		expected: []byte{4, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 'a', 'b', 'c', 2, 0, 0, 0, 0x00, 0x04},
	},
}

func TestToBytesConversion(t *testing.T) {
	for name, tc := range convertToBytesTests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			bc.WriteBytesTo(buf, tc.input)
			bs := buf.Bytes()
			if !slices.Equal(bs, tc.expected) {
				t.Errorf("%s: got %v, expected %v", name, bs, tc.expected)
			}
		})
	}
}
