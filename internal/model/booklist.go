package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/tp86/legimi-go/internal/api/protocol/encoding"
)

type BookListRequest struct {
	session
	NextPage string
	BookId   uint64
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

func (l BookListRequest) DebugFormat() string {
	var b strings.Builder
	b.WriteString("List books request:\n")
	filters := makeFilters(l)
	b.WriteString("requested page: ")
	b.WriteString(l.NextPage)
	b.WriteString("\n")
	b.WriteString("filters: ")
	for _, filter := range filters {
		b.WriteString(fmt.Sprintf("%02X ", filter.Type))
		b.WriteString(fmt.Sprintf("%02X ", filter.Subtype))
		b.WriteString(fmt.Sprintf("%02X ", filter.Data))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("requested book id: %d", l.BookId))
	return b.String()
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
	return encoding.U8Length /*filter type*/ +
		encoding.U16Length /*filter subtype*/ +
		encoding.U16Length /*filter data length*/ + encoding.EncodedLength(f.Data) /*filter data*/
}

const (
	filterBookList                = 2
	filterBookFormat              = 14
	filterBookFormatKindle uint16 = 8
	filterBookId                  = 10
	filterPagingType              = 4
	filterPagingSubtype           = 600
	filterPageSizeKey             = 3
	filterPageDefaultSize  uint32 = 500
	filterPageNextKey             = 4
)

func makePagingFilter(l BookListRequest) filter {
	data := encoding.Map{filterPageSizeKey: filterPageDefaultSize}
	if l.NextPage != "" {
		data[filterPageNextKey] = l.NextPage
	}
	return filter{filterPagingType, filterPagingSubtype, data}
}

func makeBookFilters(l BookListRequest) []filter {
	filters := []filter{
		{filterBookList, filterBookFormat, filterBookFormatKindle},
	}
	if l.BookId != 0 {
		filters = append(filters, filter{filterBookList, filterBookId, l.BookId})
	}
	return filters
}

func makeFilters(l BookListRequest) []filter {
	filters := makeBookFilters(l)
	filters = append(filters, makePagingFilter(l))
	return filters
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
	bytesRead, err := checkBookListMetadataIsSupported(r)
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

func checkBookListMetadataIsSupported(r io.Reader) (int, error) {
	var metadataType uint8
	bytesRead, err := encoding.Decode(r, &metadataType)
	if err != nil {
		return bytesRead, err
	}
	if metadataType != booklistMetadataSupportedType {
		return bytesRead, fmt.Errorf("book metadata type should be %d, found %d", booklistMetadataSupportedType, metadataType)
	}
	return bytesRead, nil
}

const (
	booklistMetadataSupportedType = 7
)

type booklistMetadata encoding.Map

var toSkipInBooklistMetadata = []int{
	encoding.U64Length,
	encoding.EmptyLength,
	encoding.U32Length,
	encoding.U64Length,
	encoding.EmptyLength,
}

func (md booklistMetadata) Decode(r io.Reader) (int, error) {
	bytesRead, err := encoding.SkipDecodeMany(r, toSkipInBooklistMetadata)
	if err != nil {
		return bytesRead, err
	}
	n, err := encoding.Decode(r, encoding.Map(md))
	bytesRead += n
	return bytesRead, err
}
