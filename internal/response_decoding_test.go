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
		t.Errorf("register response decoding: expected %v, got %v", expected, registerResponse)
	}
}

func TestSessionResponseDecoding(t *testing.T) {
	input := []byte{
		0x11, 0x00, 0x00, 0x00, 0x02, 0x40, 0x28, 0x00, 0x00, 0x00, 0x01, 0x00, 0x07, 0x00, 0x20, 0x00,
		0x00, 0x00, '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v',
	}
	expected := response.Session{Id: "1234567890abcdefghijklmnopqrstuv"}
	buf := bytes.NewBuffer(input)
	var sessionResponse response.Session
	err := packet.Decode(buf, &sessionResponse)
	if err != nil {
		t.Fatalf("decoding error: %v", err)
	}
	if sessionResponse != expected {
		t.Errorf("session response decoding: expected %v, got %v", expected, sessionResponse)
	}
}
