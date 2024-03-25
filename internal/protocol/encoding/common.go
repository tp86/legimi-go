package encoding

const (
	U8Length  = 1
	U16Length = 2
	U32Length = 4
	U64Length = 8
)

type Key = uint16

type Typed interface {
	Type() uint16
}
