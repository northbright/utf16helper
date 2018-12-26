package utf16helper_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/northbright/utf16helper"
)

func ExampleDetectUTF16BOM() {
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
		order, err := utf16helper.DetectUTF16BOM(buf)
		if err != nil {
			log.Printf("DtectUTF16BOM() error: %v", err)
			return
		}
		log.Printf("order: %v", order)
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
		if err != nil && err != utf16helper.ErrNoUTF16BOM {
			log.Printf("WriteUTF16BOM() error: %v", err)
			return
		}

		if err != utf16helper.ErrNoUTF16BOM {
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

	// Case 1: Use UTF8 bytes directly.
	// "Hello, 世界!" in UTF8 with BOM(written in Notepad.exe on Windows).
	utf8Bytes := []byte{0xEF, 0xBB, 0xBF, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0x21}
	buf := []byte{}

	// Create a bytes reader(io.Reader) from UTF8 bytes.
	reader := bytes.NewReader(utf8Bytes)
	// Create a buffer writer(io.Writer) to contain converted UTF16 bytes.
	writer := bytes.NewBuffer(buf)

	if err := utf16helper.UTF8ToUTF16Ctx(ctx, reader, writer); err != nil {
		log.Printf("UTF8ToUTF16Ctx() error: %v", err)
		return
	}

	log.Printf("%X", writer.Bytes())

	// Case 2: Convert string to UTF8 bytes.
	utf8Bytes = []byte(`Hello, 世界!`)
	buf = []byte{}

	// Create a bytes reader(io.Reader) from UTF8 bytes.
	reader = bytes.NewReader(utf8Bytes)
	// Create a buffer writer(io.Writer) to contain converted UTF16 bytes.
	writer = bytes.NewBuffer(buf)

	if err := utf16helper.UTF8ToUTF16Ctx(ctx, reader, writer); err != nil {
		log.Printf("UTF8ToUTF16Ctx() error: %v", err)
		return
	}

	log.Printf("%X", writer.Bytes())

	// Output:
}

func ExampleUTF16ToUTF8Ctx() {
	// Get an empty context.
	ctx := context.Background()
	// "Hello, 世界!" in UTF16 with UTF16LE(written in Notepad.exe on Windows).
	utf16Bytes := []byte{0xFF, 0xFE, 0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00,
		0x6f, 0x00, 0x2c, 0x00, 0x20, 0x00, 0x16, 0x4e, 0x4c, 0x75, 0x21, 0x00}
	buf := []byte{}

	// Create a bytes reader(io.Reader) from UTF16 bytes.
	reader := bytes.NewReader(utf16Bytes)
	// Create a buffer writer(io.Writer) to contain converted UTF8 bytes.
	writer := bytes.NewBuffer(buf)

	if err := utf16helper.UTF16ToUTF8Ctx(ctx, reader, writer, false); err != nil {
		log.Printf("UTF16ToUTF8Ctx() error: %v", err)
		return
	}

	log.Printf("%X", writer.Bytes())
	// Output string converted from UTF8 bytes.
	log.Printf("string: %s", writer.String())
	// Output:
}
