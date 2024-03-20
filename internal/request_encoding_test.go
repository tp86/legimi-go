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
		t.Errorf("register request encoding: expected %v, got %v", expected, b)
	}
}

func TestSessionRequestEncoding(t *testing.T) {
	input := struct {
		login, password string
		kindleId        uint64
	}{
		"login", "password", 12345678,
	}
	expected := map[int][]byte{
		-1: {0x11, 0x00, 0x00, 0x00, 0x50, 0x00, 0x54, 0x00, 0x00, 0x00, 0x06, 0x00},
		0:  {0x05, 0x00, 0x00, 0x00, 'l', 'o', 'g', 'i', 'n'},
		1:  {0x08, 0x00, 0x00, 0x00, 'p', 'a', 's', 's', 'w', 'o', 'r', 'd'},
		2:  {0x08, 0x00, 0x00, 0x00, 0x4e, 0x61, 0xbc, 0x00, 0x00, 0x00, 0x00, 0x00},
		3:  {0x0d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x38, 0x2e, 0x35, 0x20, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73},
		4:  {0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		5:  {0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	buf := new(bytes.Buffer)
	sessionRequest := request.NewSessionRequest(input.login, input.password, input.kindleId)
	err := packet.Encode(buf, sessionRequest)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	b := buf.Bytes()
	expectedHeader := expected[-1]
	bytesRead := len(expectedHeader)
	header := b[:len(expectedHeader)]
	if !slices.Equal(header, expectedHeader) {
		t.Errorf("session request header encoding: expected %v, got %v", expectedHeader, header)
	}
	for bytesRead < buf.Len() {
		keyBytes := b[bytesRead : bytesRead+2]
		key := int(keyBytes[0]) + int(keyBytes[1])<<8
		bytesRead += 2
		lengthBytes := b[bytesRead : bytesRead+4]
		length := int(lengthBytes[0]) + int(lengthBytes[1])<<8 + int(lengthBytes[2])<<16 + int(lengthBytes[3])<<24
		valueBytes := b[bytesRead : bytesRead+4+length]
		expectedValue := expected[key]
		if !slices.Equal(valueBytes, expectedValue) {
			t.Errorf("session request key-value pair encoding: key %d, expected %v, got %v", key, expectedValue, valueBytes)
		}
		bytesRead += 4 + length
	}
}
