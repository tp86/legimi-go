package byteconverter

import "io"

type Byte struct {
	Value uint8
}

func (b *Byte) WriteBytesTo(w io.ByteWriter) {
	w.WriteByte(b.Value)
}

func (b *Byte) ReadBytesFrom(r io.ByteReader) error {
	rb, err := r.ReadByte()
	if err != nil {
		return err
	}
	b.Value = rb
	return nil
}
