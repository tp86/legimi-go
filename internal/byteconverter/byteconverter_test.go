package byteconverter_test

import (
	"bytes"
	"slices"
	"testing"

	bc "github.com/tp86/legimi-go/internal/byteconverter"
)

func TestByteWriting(t *testing.T) {
	buf := new(bytes.Buffer)
	b := &bc.Byte{1}
	b.WriteBytesTo(buf)
	if !slices.Equal(buf.Bytes(), []byte{1}) {
		t.Errorf("byte: expected: %v, got: %v", 1, 1)
	}
}

func TestByteReading(t *testing.T) {
}
