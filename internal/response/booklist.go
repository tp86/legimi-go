package response

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type BookList protocol.Array[BookMetadata]

type BookMetadata struct {
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
		blSlice = make([]BookMetadata, 0)
	}
	array := protocol.Array[BookMetadata](blSlice)
	n, err := protocol.Decode(r, &array)
	if err != nil {
		return n, err
	}
	*bl = BookList(array)
	return n, err
}

func (bm *BookMetadata) Decode(r io.Reader) (int, error) {
	var bytesRead int
	// skip type of item, only 7 supported
	// TODO check type of item
	bytesRead, err := protocol.SkipDecode(r, protocol.U8Length)
	if err != nil {
		return bytesRead, err
	}
	n, err := protocol.Decode(r, protocol.WithLength{Value: metadata{
		10: &bm.Id,
		11: &bm.Title,
		0:  &bm.Author,
		13: &bm.Version,
		30: &bm.Downloaded,
		34: &bm.NextPage,
	}})
	bytesRead += n
	return bytesRead, err
}

type metadata protocol.Map

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

func (md metadata) Decode(r io.Reader) (int, error) {
	var bytesRead int
	for _, skip := range toSkip {
		n, err := protocol.SkipDecode(r, skip)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
	}
	m := protocol.Map(md)
	n, err := protocol.Decode(r, m)
	bytesRead += n
	return bytesRead, err
}
