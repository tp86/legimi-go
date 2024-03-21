package response

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type BookList []BookDetails

type BookDetails struct {
	Id            uint64
	Title         string
	Author        string
	Version       uint64
	Downloaded    bool
	NextPageToken string
}

func (bl *BookList) Decode(r io.Reader) (int, error) {
	var decoders = make([]BookDetails, 0)
	n, err := protocol.Decode(r, &decoders)
	if err != nil {
		return n, err
	}
	// bookDetails := make([]BookDetails, len(decoders))
	// for i, decoder := range decoders {
	//	bookDetails[i] = *decoder.(*BookDetails)
	// }
	*bl = BookList(decoders)
	return n, err
}

func (bl BookList) Type() uint16 {
	return 0x001C
}

func (bd *BookDetails) Decode(r io.Reader) (int, error) {
	return 0, nil
}
