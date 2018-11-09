package utf16helper_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func ExampleWriteUTF16BOM() {
	orders := []binary.ByteOrder{
		binary.LittleEndian,
		binary.BigEndian,
		// No order
		nil,
	}

	for _, order := range orders {
		buf := &bytes.Buffer{}

		err := utf16helper.WriteUTF16BOM(order, buf)
		if err != nil && err != utf16helper.ErrNoBOM {
			log.Printf("WriteUTF16BOM() error: %v", err)
			return
		}

		if err != utf16helper.ErrNoBOM {
			fmt.Printf("written BOM: 0x%X\n", buf.Bytes())
		} else {
			fmt.Printf("no BOM\n")
		}

	}

	// Output:
	// written BOM: 0xFFFE
	// written BOM: 0xFEFF
	// no BOM
}
