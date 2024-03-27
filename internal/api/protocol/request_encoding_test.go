package protocol_test

import (
	"bytes"
	"slices"
	"testing"

	"github.com/tp86/legimi-go/internal/api/protocol"
	"github.com/tp86/legimi-go/internal/model"
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

	req := model.NewRegisterRequest(input.login, input.password, input.serial)
	encodeAndCheck(t, "register", req, expected)
}

func TestGetSessionRequestEncoding(t *testing.T) {
	input := struct {
		login, password string
		kindleId        uint64
	}{"login", "password", 12345678}
	expected := map[int][]byte{
		encodedMapHeaderKey: {0x11, 0x00, 0x00, 0x00, 0x50, 0x00, 0x54, 0x00, 0x00, 0x00, 0x06, 0x00},
		0:                   {0x05, 0x00, 0x00, 0x00, 'l', 'o', 'g', 'i', 'n'},
		1:                   {0x08, 0x00, 0x00, 0x00, 'p', 'a', 's', 's', 'w', 'o', 'r', 'd'},
		2:                   {0x08, 0x00, 0x00, 0x00, 0x4e, 0x61, 0xbc, 0x00, 0x00, 0x00, 0x00, 0x00},
		3:                   {0x0d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x38, 0x2e, 0x35, 0x20, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73},
		4:                   {0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		5:                   {0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	req := model.NewGetSessionRequest(input.login, input.password, input.kindleId)
	encodeAndCheckMap(t, "session", req, expected)
}

func TestBookListRequestEncoding(t *testing.T) {
	input := struct {
		sessionId string
		nextPage  string
	}{sessionId: "1234567890abcdefghijklmnopqrstuv"}
	expected := map[int][]byte{
		encodedMapHeaderKey: {0x11, 0x00, 0x00, 0x00, 0x1a, 0x00, 0x39, 0x00, 0x00, 0x00, 0x02, '1', '2', '3', '4', '5',
			'6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
			'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 0x02, 0x0e, 0x00, 0x02, 0x00,
			0x08, 0x00, 0x04, 0x58, 0x02, 0x0c, 0x00, 0x01, 0x00},
		3: {0x04, 0x00, 0x00, 0x00, 0xf4, 0x01, 0x00, 0x00},
	}

	req := model.NewBookListRequest(input.sessionId)
	encodeAndCheckMap(t, "list books", req, expected)

	addNewFilter := func(expected map[int][]byte, key int, filter []byte) {
		nextPageFilterLength := 2 /*key*/ + uint8(len(filter)) /*value*/
		expectedHeader := expected[encodedMapHeaderKey]
		// update length of payload
		expectedHeader[6] += nextPageFilterLength
		// update length of filters
		expectedHeader[len(expectedHeader)-4] += nextPageFilterLength
		// update count of map entries
		expectedHeader[len(expectedHeader)-2] += 1
		// add map entry
		expected[key] = filter
	}
	nextPageFilter := []byte{0x08, 0x00, 0x00, 0x00, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}
	addNewFilter(expected, 4, nextPageFilter)

	req.NextPage = "12345678"
	encodeAndCheckMap(t, "list books", req, expected)
}

func TestBookDownloadDetailsRequestEncoding(t *testing.T) {
	input := struct {
		sessionId   string
		bookId      uint64
		bookVersion uint64
	}{
		sessionId:   "1234567890abcdefghijklmnopqrstuv",
		bookId:      12345678,
		bookVersion: 4,
	}
	expected := []byte{
		0x11, 0x00, 0x00, 0x00, 0xc8, 0x00, 0x40, 0x00, 0x00, 0x00, 0x4e, 0x61, 0xbc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, '1', '2', '3', '4', '5', '6',
		'7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
		'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	req := model.NewBookDownloadDetailsRequest(input.sessionId, input.bookId, input.bookVersion)
	encodeAndCheck(t, "book download details", req, expected)
}

func encodeAndCheck(t *testing.T, name string, req protocol.Request, expected []byte) {
	buf := new(bytes.Buffer)
	err := protocol.Encode(buf, req)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	b := buf.Bytes()
	if !slices.Equal(b, expected) {
		t.Errorf("%s request encoding: expected %v, got %v", name, expected, b)
	}
}

func encodeAndCheckMap(t *testing.T, name string, req protocol.Request, expected map[int][]byte) {
	buf := new(bytes.Buffer)
	err := protocol.Encode(buf, req)
	if err != nil {
		t.Fatalf("encoding error: %v", err)
	}
	checkMapEqual(t, name, expected, buf.Bytes())
}

const (
	encodedMapHeaderKey = -1
)

func checkMapEqual(t *testing.T, name string, expected map[int][]byte, actual []byte) {
	expectedHeader := expected[encodedMapHeaderKey]
	bytesRead := len(expectedHeader)
	header := actual[:bytesRead]
	if !slices.Equal(header, expectedHeader) {
		t.Errorf("%s request header encoding: expected %v, got %v", name, expectedHeader, header)
	}
	for bytesRead < len(actual) {
		keyBytes := actual[bytesRead : bytesRead+2]
		key := int(keyBytes[0]) + int(keyBytes[1])<<8
		bytesRead += 2
		lengthBytes := actual[bytesRead : bytesRead+4]
		length := int(lengthBytes[0]) + int(lengthBytes[1])<<8 + int(lengthBytes[2])<<16 + int(lengthBytes[3])<<24
		valueBytes := actual[bytesRead : bytesRead+4+length]
		expectedValue := expected[key]
		if !slices.Equal(valueBytes, expectedValue) {
			t.Errorf("%s request key-value pair encoding: key %d, expected %v, got %v", name, key, expectedValue, valueBytes)
		}
		bytesRead += 4 + length
	}
}
