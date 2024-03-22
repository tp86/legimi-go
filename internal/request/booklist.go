package request

import (
	"io"

	"github.com/tp86/legimi-go/internal/protocol"
)

type BookList struct {
	SessionId string
	NextPage  string
}

func (l BookList) Encode(w io.Writer) error {
	filters := makeFilters(l)
	for _, value := range []any{
		uint8(len(filters)),
		l.SessionId,
	} {
		err := protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	for _, filter := range filters {
		err := protocol.Encode(w, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l BookList) EncodedLength() int {
	var filtersLength int
	for _, filter := range makeFilters(l) {
		filtersLength += protocol.EncodedLength(filter)
	}
	return protocol.U8Length +
		protocol.EncodedLength(l.SessionId) +
		filtersLength
}

func (l BookList) Type() uint16 {
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
		uint16(protocol.EncodedLength(f.Data)),
		f.Data,
	} {
		err := protocol.Encode(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f filter) EncodedLength() int {
	return protocol.U8Length +
		protocol.U16Length +
		protocol.U16Length + protocol.EncodedLength(f.Data)
}

func makeFilters(l BookList) []filter {
	data := protocol.Map{3: uint32(500)}
	if l.NextPage != "" {
		data[4] = l.NextPage
	}
	secondFilter := filter{4, 600, data}
	return []filter{
		{2, 14, uint16(8)},
		secondFilter,
	}
}
