package response

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type BookList protocol.Array[BookDetails]

type BookDetails struct {
	Id         uint64
	Title      string
	Author     string
	Version    uint64
	Downloaded bool
	NextPage   string
}

func (bl BookList) Type() uint16 {
	return 0x001C
}

func (bl *BookList) Decode(r io.Reader) (int, error) {
	blSlice := *bl
	if blSlice == nil {
		blSlice = make([]BookDetails, 0)
	}
	array := protocol.Array[BookDetails](blSlice)
	n, err := protocol.Decode(r, &array)
	if err != nil {
		return n, err
	}
	*bl = BookList(array)
	return n, err
}

func (bd *BookDetails) Decode(r io.Reader) (int, error) {
	var bytesRead int
	// skip type of item, only 7 supported
	// TODO check type of item
	bytesRead, err := protocol.SkipDecode(r, protocol.U8Length)
	if err != nil {
		return bytesRead, err
	}
	n, err := protocol.Decode(r, protocol.WithLength{Value: details{
		10: &bd.Id,
		11: &bd.Title,
		0:  &bd.Author,
		13: &bd.Version,
		30: &bd.Downloaded,
		34: &bd.NextPage,
	}})
	bytesRead += n
	return bytesRead, err
}

type details protocol.Map

var (
	emptyLength = protocol.WithLength{}.EncodedLength()
	toSkip      = []int{
		protocol.U64Length,
		emptyLength,
		protocol.U32Length,
		protocol.U64Length,
		emptyLength,
	}
)

func (d details) Decode(r io.Reader) (int, error) {
	var bytesRead int
	for _, value := range toSkip {
		n, err := protocol.SkipDecode(r, value)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
	}
	m := protocol.Map(d)
	n, err := protocol.Decode(r, m)
	bytesRead += n
	return bytesRead, err
}
