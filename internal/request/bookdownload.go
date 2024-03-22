package request

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type BookDownloadDetails struct {
	session
	bookId  uint64
	version uint64
}

func NewBookDownloadDetailsRequest(sessionId string, bookId uint64, bookVersion uint64) BookDownloadDetails {
	return BookDownloadDetails{
		session: session{id: sessionId},
		bookId:  bookId,
		version: bookVersion,
	}
}

func (bd BookDownloadDetails) Encode(w io.Writer) error {
	for _, value := range bd.data() {
		err := protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bd BookDownloadDetails) EncodedLength() int {
	var totalLength int
	for _, value := range bd.data() {
		totalLength += protocol.EncodedLength(value)
	}
	return totalLength
}

func (bd BookDownloadDetails) Type() uint16 {
	return 0x00C8
}

func (bd BookDownloadDetails) data() []any {
	return []any{
		bd.bookId,
		bd.version,
		bd.session.id,
		false,
		false,
		uint64(0xFFFFFFFFFFFFFFFF),
		// XXX
		uint16(1),
		uint16(0),
		uint16(0),
	}
}
