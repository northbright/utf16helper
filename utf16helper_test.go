package utf16helper_test

import (
	"bytes"
	"log"

	"github.com/northbright/utf16helper"
)

func ExampleReadUTF16BOM() {
	arr := [][]byte{
		// UTF-16 string: "服" with little endian BOM.
		[]byte{0xFF, 0xFE, 0x0d, 0x67},
		// UTF-16 string: "务" with big endian BOM.
		[]byte{0xFE, 0xFF, 0x52, 0xa1},
		// UTF-16 string: "器" without BOM.
		[]byte{0x68, 0x56},
	}

	for _, v := range arr {
		buf := bytes.NewBuffer(v)
		bom, order, err := utf16helper.ReadUTF16BOM(buf)
		if err != nil {
			log.Printf("ReadUTFBOM() error: %v", err)
			return
		}
		log.Printf("BOM(first 2 bytes): %X, order: %v", bom, order)
	}

	// Output:
}
