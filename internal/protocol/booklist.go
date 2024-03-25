package protocol

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol/encoding"
)

type BookListRequest struct {
	session
	NextPage string
}

func NewBookListRequest(sessionId string) BookListRequest {
	return BookListRequest{
		session: session{id: sessionId},
	}
}

func (l BookListRequest) Encode(w io.Writer) error {
	filters := makeFilters(l)
	for _, value := range []any{
		uint8(len(filters)),
		l.session.id,
	} {
		err := encoding.Encode(w, value)
		if err != nil {
			return err
		}
	}
	for _, filter := range filters {
		err := encoding.Encode(w, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l BookListRequest) EncodedLength() int {
	var filtersLength int
	for _, filter := range makeFilters(l) {
		filtersLength += encoding.EncodedLength(filter)
	}
	return encoding.U8Length +
		encoding.EncodedLength(l.session.id) +
		filtersLength
}

func (l BookListRequest) Type() uint16 {
	return 0x001a
}

type filter struct {
	Type    uint8
	Subtype uint16
	Data    any
}

func (f filter) Encode(w io.Writer) error {
	for _, value := range []any{
		f.Type,
		f.Subtype,
		uint16(encoding.EncodedLength(f.Data)),
		f.Data,
	} {
		err := encoding.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f filter) EncodedLength() int {
	return encoding.U8Length +
		encoding.U16Length +
		encoding.U16Length + encoding.EncodedLength(f.Data)
}

func makeFilters(l BookListRequest) []filter {
	data := encoding.Map{3: uint32(500)}
	if l.NextPage != "" {
		data[4] = l.NextPage
	}
	secondFilter := filter{4, 600, data}
	return []filter{
		{2, 14, uint16(8)},
		secondFilter,
	}
}

type BookList encoding.Array[BookMetadata]

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
	array := encoding.Array[BookMetadata](blSlice)
	n, err := encoding.Decode(r, &array)
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
	bytesRead, err := encoding.SkipDecode(r, encoding.U8Length)
	if err != nil {
		return bytesRead, err
	}
	n, err := encoding.Decode(r, encoding.WithLength{Value: booklistMetadata{
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

type booklistMetadata encoding.Map

var toSkipInBooklistMetadata = []int{
	encoding.U64Length,
	emptyLength,
	encoding.U32Length,
	encoding.U64Length,
	emptyLength,
}

func (md booklistMetadata) Decode(r io.Reader) (int, error) {
	var bytesRead int
	for _, skip := range toSkipInBooklistMetadata {
		n, err := encoding.SkipDecode(r, skip)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}
	}
	m := encoding.Map(md)
	n, err := encoding.Decode(r, m)
	bytesRead += n
	return bytesRead, err
}
