package packet

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
	"github.com/tp86/legimi-go/internal/request"
	"github.com/tp86/legimi-go/internal/response"
)

func Encode(w io.Writer, req request.Request) error {
	for _, value := range []any{
		uint32(17),
		req.Type(),
		protocol.WithLength{Value: req},
	} {
		err := protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Decode(r io.Reader, resp response.Response) error {
	// TODO handle errors
	protocol.SkipDecode(r, protocol.U32Length)
	// TODO check protocol version?
	var responseType uint16
	protocol.Decode(r, &responseType)
	// TODO handle error response type
	if responseType != resp.Type() {
		return fmt.Errorf("unexpected response type: %d, expected: %d", responseType, resp.Type())
	}
	var packetLength uint32
	protocol.Decode(r, &packetLength)
	bytesRead, err := protocol.Decode(r, resp)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	if bytesRead != int(packetLength) {
		return fmt.Errorf("bytes read: %d, expected: %d", bytesRead, packetLength)
	}
	return nil
}
