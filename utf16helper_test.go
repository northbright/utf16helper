package utf16helper_test

import (
	"bytes"
	"context"
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

func ExampleWriteUTF8BOM() {
	buf := &bytes.Buffer{}

	err := utf16helper.WriteUTF8BOM(buf)
	if err != nil {
		log.Printf("WriteUTF8BOM() error: %v", err)
		return
	}

	fmt.Printf("written BOM: 0x%X\n", buf.Bytes())

	// Output:
	// written BOM: 0xEFBBBF
}

func ExampleUTF8ToUTF16Ctx() {
	// Get an empty context.
	ctx := context.Background()
	// "Hello, 世界!" in UTF8 with BOM(written in Notepad.exe on Windows).
	utf8Bytes := []byte{0xEF, 0xBB, 0xBF, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0x21}
	buf := []byte{}

	reader := bytes.NewReader(utf8Bytes)
	writer := bytes.NewBuffer(buf)

	if err := utf16helper.UTF8ToUTF16Ctx(ctx, reader, writer); err != nil {
		log.Printf("UTF8ToUTF16Ctx() error: %v", err)
		return
	}

	log.Printf("%X\n", writer.Bytes())
	// Output:
}
