package protocol

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol/encoding"
)

func Encode(w io.Writer, req Request) error {
	for _, value := range []any{
		ProtocolVersion,
		req.Type(),
		encoding.WithLength{Value: req},
	} {
		err := encoding.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Decode(r io.Reader, resp Response) error {
	_, err := checkPacketIsSupported(r)
	if err != nil {
		return fmt.Errorf("received packet not supported: %v", err)
	}

	var responseType uint16
	encoding.Decode(r, &responseType)
	if responseType != resp.Type() {
		if errorResponse, isErrorResponse := errorResponses[responseType]; isErrorResponse {
			return fmt.Errorf("received error response: %v", errorResponse)
		}
		return fmt.Errorf("unexpected response type: %d, expected: %d", responseType, resp.Type())
	}

	var packetLength uint32
	encoding.Decode(r, &packetLength)
	if packetLength == 0 {
		return fmt.Errorf("empty packet body")
	}

	bytesRead, err := encoding.Decode(r, resp)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	if bytesRead != int(packetLength) {
		return fmt.Errorf("bytes read: %d, expected: %d", bytesRead, packetLength)
	}
	return nil
}

func checkPacketIsSupported(r io.Reader) (int, error) {
	var protocolVersion uint32
	bytesRead, err := encoding.Decode(r, &protocolVersion)
	if err != nil {
		return bytesRead, err
	}
	if protocolVersion != ProtocolVersion {
		return bytesRead, fmt.Errorf("protocol version should be %d, received %d", ProtocolVersion, protocolVersion)
	}
	return bytesRead, nil
}
