package byteconverter_

import "io"

type ByteReaderFrom interface {
	readBytesFrom(r io.ByteReader)
}

func ReadBytesFrom[V numbers | string](r io.ByteReader, c func(v V) ByteReaderFrom) V {
	b, _ := r.ReadByte()
	println(b)
	return V(b)
}
