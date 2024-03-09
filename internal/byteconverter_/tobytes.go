package byteconverter_

import "io"

type ByteWriterTo interface {
	writeBytesTo(w io.ByteWriter)
}

func WriteBytesTo(w io.ByteWriter, val ByteWriterTo) {
	val.writeBytesTo(w)
}
