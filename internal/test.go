package main

import (
	"bytes"
	"fmt"

	"github.com/tp86/legimi-go/internal/packet"
	"github.com/tp86/legimi-go/internal/request"
)

func printBytes(bs []byte) {
	for i, b := range bs {
		fmt.Printf("%02X ", b)
		if i%16 == 15 {
			fmt.Println()
		}
	}
}

func main() {
	req := request.Registration{
		Login:          "",
		Password:       ``,
		KindleSerialNo: "",
	}
	packet := packet.NewPacket(req)

	buf := new(bytes.Buffer)
	packet.WriteBytesTo(buf)
	printBytes(buf.Bytes())
}
