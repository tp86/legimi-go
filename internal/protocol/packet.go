package protocol

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol/encoding"
)

func Encode(w io.Writer, req Request) error {
	for _, value := range []any{
		Version,
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
	// TODO handle errors
	encoding.SkipDecode(r, encoding.U32Length)
	// TODO check protocol version?
	var responseType uint16
	encoding.Decode(r, &responseType)
	// TODO handle error response type
	if responseType != resp.Type() {
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
