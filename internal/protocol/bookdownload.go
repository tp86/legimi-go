package protocol

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol/encoding"
)

type BookDownloadDetailsRequest struct {
	session
	bookId  uint64
	version uint64
}

func NewBookDownloadDetailsRequest(sessionId string, bookId uint64, bookVersion uint64) BookDownloadDetailsRequest {
	return BookDownloadDetailsRequest{
		session: session{id: sessionId},
		bookId:  bookId,
		version: bookVersion,
	}
}

func (bd BookDownloadDetailsRequest) Encode(w io.Writer) error {
	for _, value := range bd.data() {
		err := encoding.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bd BookDownloadDetailsRequest) EncodedLength() int {
	var totalLength int
	for _, value := range bd.data() {
		totalLength += encoding.EncodedLength(value)
	}
	return totalLength
}

func (bd BookDownloadDetailsRequest) Type() uint16 {
	return 0x00C8
}

func (bd BookDownloadDetailsRequest) data() []any {
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

type BookDownloadDetails struct {
	Url  string
	Size uint64
}

// TODO refactor

const (
	detailsCountSupported = 1
	typeSupported         = 4
)

func (bdd *BookDownloadDetails) Decode(r io.Reader) (int, error) {
	var bytesRead int
	for _, skip := range toSkipInBookDownloadDetailsHeader {
		n, err := encoding.SkipDecode(r, skip)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
	}
	var count uint32
	n, err := encoding.Decode(r, &count)
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}
	if count != detailsCountSupported {
		return bytesRead, fmt.Errorf("there should be only %d download details in response, received %d", detailsCountSupported, count)
	}
	var typ uint8
	n, err = encoding.Decode(r, &typ)
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}
	if typ != typeSupported {
		return bytesRead, fmt.Errorf("download details type should be %d, found %d", typeSupported, typ)
	}
	n, err = encoding.Decode(r, encoding.Map{
		0: &bdd.Url,
		2: &bdd.Size,
	})
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}
	for _, skip := range toSkipInBookDownloadDetailsFooter {
		n, err := encoding.SkipDecode(r, skip)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
	}
	return bytesRead, err
}

func (bdd BookDownloadDetails) Type() uint16 {
	return 0x0018
}

var (
	toSkipInBookDownloadDetailsHeader = []int{encoding.U8Length, encoding.U32Length}
	toSkipInBookDownloadDetailsFooter = []int{
		emptyLength,
		encoding.U64Length,
		encoding.U32Length,
	}
)
