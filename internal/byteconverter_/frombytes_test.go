package byteconverter__test

import (
	"bytes"
	"testing"

	bc "github.com/tp86/legimi-go/internal/byteconverter"
)

var convertFromBytesTests = map[string]struct {
	input    []byte
	expected any
}{
	"byte":    {[]byte{1, 0, 0, 0, 1}, byte(1)},
	"integer": {[]byte{1, 0, 0, 0, 2, 0, 0, 0}, uint32(2)},
}

func TestFromBytesConversion(t *testing.T) {
	for name, tc := range convertFromBytesTests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(tc.input)
			output := bc.ReadBytesFrom(buf, bc.Byte)
			if output != tc.expected {
				t.Errorf("%s: got %v, expected %v", name, output, tc.expected)
			}
		})
	}
}
