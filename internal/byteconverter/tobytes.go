package byteconverter

import "io"

type ByteWriterTo interface {
	writeBytesTo(w io.ByteWriter)
}

func WriteBytesTo(w io.ByteWriter, val ByteWriterTo) {
	val.writeBytesTo(w)
}

func Byte(b uint8) ByteWriterTo {
	return bbyte{val: b}
}

func Short(s uint16) ByteWriterTo {
	return bshort{val: s}
}

func Int(i uint32) ByteWriterTo {
	return bint{val: i}
}

func Long(l uint64) ByteWriterTo {
	return blong{val: l}
}
func String(s string) ByteWriterTo {
	return bstring[length]{val: s}
}

func ShortString(s string) ByteWriterTo {
	return bstring[shortLength]{val: s}
}

func Sequence(bvals ...ByteWriterTo) ByteWriterTo {
	return sequence(bvals)
}
