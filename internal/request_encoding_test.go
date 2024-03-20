package internal_test

import (
	"bytes"
	"slices"
	"testing"

	"github.com/tp86/legimi-go/internal/packet"
	"github.com/tp86/legimi-go/internal/request"
)

func TestRegisterRequestEncoding(t *testing.T) {
	input := struct {
		login, password, serial string
	}{"login", "password", "12345678"}
	expected := []byte{
		0x11, 0x00, 0x00, 0x00, 0x42, 0x00, 0x3D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 'l', 'o', 'g', 'i', 'n', 0x08, 0x00, 'p',
		'a', 's', 's', 'w', 'o', 'r', 'd', 0x20, 0x00, 'K', 'i', 'n', 'd', 'l', 'e', '|',
		'|', 'K', 'i', 'n', 'd', 'l', 'e', '|', '|', '1', '2', '3', '4', '5', '6', '7',
		'8', '|', '|', 'K', 'i', 'n', 'd', 'l', 'e', 0x00, 0x00,
	}
	buf := new(bytes.Buffer)
	regReq := request.NewRegisterRequest(input.login, input.password, input.serial)
	err := packet.Encode(buf, regReq)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	b := buf.Bytes()
	if !slices.Equal(b, expected) {
		t.Errorf("Register request encoding: expected %v, got %v", expected, b)
	}
}
