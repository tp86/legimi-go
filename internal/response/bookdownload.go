package response

import (
	"fmt"
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

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
		n, err := protocol.SkipDecode(r, skip)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
	}
	var count uint32
	n, err := protocol.Decode(r, &count)
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}
	if count != detailsCountSupported {
		return bytesRead, fmt.Errorf("there should be only %d download details in response, received %d", detailsCountSupported, count)
	}
	var typ uint8
	n, err = protocol.Decode(r, &typ)
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}
	if typ != typeSupported {
		return bytesRead, fmt.Errorf("download details type should be %d, found %d", typeSupported, typ)
	}
	n, err = protocol.Decode(r, protocol.Map{
		0: &bdd.Url,
		2: &bdd.Size,
	})
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}
	for _, skip := range toSkipInBookDownloadDetailsFooter {
		n, err := protocol.SkipDecode(r, skip)
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
	toSkipInBookDownloadDetailsHeader = []int{protocol.U8Length, protocol.U32Length}
	toSkipInBookDownloadDetailsFooter = []int{
		emptyLength,
		protocol.U64Length,
		protocol.U32Length,
	}
)
