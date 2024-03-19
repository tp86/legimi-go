package debug

import "fmt"

func PrintBytes(bs []byte) {
	for i, b := range bs {
		fmt.Printf("%02X ", b)
		if i%16 == 15 {
			fmt.Println()
		}
	}
	fmt.Println()
}
