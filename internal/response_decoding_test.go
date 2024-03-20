package internal_test

import (
	"bytes"
	"testing"

	"github.com/tp86/legimi-go/internal/packet"
	"github.com/tp86/legimi-go/internal/response"
)

func TestRegisterResponseDecoding(t *testing.T) {
	input := []byte{
		0x11, 0x00, 0x00, 0x00, 0x00, 0x40, 0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x06, 0x00, 0x08, 0x00,
		0x00, 0x00, 0x4e, 0x61, 0xbc, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	expected := response.Register{KindleId: 12345678}
	buf := bytes.NewBuffer(input)
	var registerResponse response.Register
	err := packet.Decode(buf, &registerResponse)
	if err != nil {
		t.Fatalf("decoding error: %v", err)
	}
	if registerResponse != expected {
		t.Errorf("Register response decoding: expected %v, got %v", expected, registerResponse)
	}
}
