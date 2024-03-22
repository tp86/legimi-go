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

func TestGetSessionRequestEncoding(t *testing.T) {
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
	sessionRequest := request.NewGetSessionRequest(input.login, input.password, input.kindleId)
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

func TestBookListRequestEncoding(t *testing.T) {
	sessionId := "1234567890abcdefghijklmnopqrstuv"
	expected := []byte{
		0x11, 0x00, 0x00, 0x00, 0x1a, 0x00, 0x39, 0x00, 0x00, 0x00, 0x02, '1', '2', '3', '4', '5',
		'6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
		'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 0x02, 0x0e, 0x00, 0x02, 0x00,
		0x08, 0x00, 0x04, 0x58, 0x02, 0x0c, 0x00, 0x01, 0x00, 0x03, 0x00, 0x04, 0x00, 0x00, 0x00, 0xf4,
		0x01, 0x00, 0x00,
	}
	buf := new(bytes.Buffer)
	listBooksRequest := request.NewBookListRequest(sessionId)
	err := packet.Encode(buf, listBooksRequest)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	b := buf.Bytes()
	if !slices.Equal(b, expected) {
		t.Errorf("list books request encoding: expected %v, got %v", expected, b)
	}
	listBooksRequest.NextPage = "12345678"
	nextPageFilter := []byte{0x04, 0x00, 0x08, 0x00, 0x00, 0x00, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}
	nextPageFilterLength := uint8(len(nextPageFilter))
	expected[6] += nextPageFilterLength
	expected[len(expected)-14] += nextPageFilterLength
	expected[len(expected)-12] += 1
	expected = append(expected, nextPageFilter...)
	buf.Reset()
	err = packet.Encode(buf, listBooksRequest)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	b = buf.Bytes()
	if !slices.Equal(b, expected) {
		t.Errorf("list books request encoding: expected %v, got %v", expected, b)
	}
}

func TestBookDownloadDetailsRequestEncoding(t *testing.T) {
	sessionId := "1234567890abcdefghijklmnopqrstuv"
	bookId := uint64(12345678)
	bookVersion := uint64(4)
	expected := []byte{
		0x11, 0x00, 0x00, 0x00, 0xc8, 0x00, 0x40, 0x00, 0x00, 0x00, 0x4e, 0x61, 0xbc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, '1', '2', '3', '4', '5', '6',
		'7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
		'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	buf := new(bytes.Buffer)
	bookDownloadDetailsRequest := request.NewBookDownloadDetailsRequest(sessionId, bookId, bookVersion)
	err := packet.Encode(buf, bookDownloadDetailsRequest)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	b := buf.Bytes()
	if !slices.Equal(b, expected) {
		t.Errorf("book download details request encoding: expected %v, got %v", expected, b)
	}
}
