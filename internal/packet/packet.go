package packet

import (
	"io"

	bc "github.com/tp86/legimi-go/internal/byteconverter"
)

type Type = uint16

const ProtocolVersion = uint32(17)

type ToBytesMapper interface {
	Type() Type
	ToBytes() []byte
}

type Packet struct {
	version    *bc.RawInt
	packetType *bc.RawShort
	payload    *bc.Bytes
}

func NewPacket(m ToBytesMapper) Packet {
	return Packet{
		version:    &bc.RawInt{Value: ProtocolVersion},
		packetType: &bc.RawShort{Value: m.Type()},
		payload:    &bc.Bytes{Value: m.ToBytes()},
	}
}

func (p Packet) WriteBytesTo(w io.ByteWriter) {
	p.version.WriteBytesTo(w)
	p.packetType.WriteBytesTo(w)
	p.payload.WriteBytesTo(w)
}
