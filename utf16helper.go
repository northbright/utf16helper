package utf16helper

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/northbright/byteorder"
)

var (
	UTF16LE       = [2]byte{0xFF, 0xFE}
	UTF16BE       = [2]byte{0xFE, 0xFF}
	UTF8BOM       = [3]byte{0xEF, 0xBB, 0xBF}
	ErrNoUTF16BOM = fmt.Errorf("No UTF-16 BOM found")
)

func DetectUTF16BOM(r io.Reader) (binary.ByteOrder, error) {
	var buf []byte

	reader := bufio.NewReader(r)

	// Read first 2 bytes.
	for i := 0; i < 2; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		buf = append(buf, b)
	}

	switch {
	case buf[0] == 0xFF && buf[1] == 0xFE:
		return binary.LittleEndian, nil
	case buf[0] == 0xFE && buf[1] == 0xFF:
		return binary.BigEndian, nil
	default:
		return nil, ErrNoUTF16BOM
	}
}

func WriteUTF16BOM(order binary.ByteOrder, dst io.Writer) error {
	var BOM []byte

	switch order {
	case nil:
		return ErrNoUTF16BOM
	case binary.LittleEndian:
		BOM = UTF16LE[0:2]
	case binary.BigEndian:
		BOM = UTF16BE[0:2]
	default:
		return ErrNoUTF16BOM
	}

	_, err := dst.Write(BOM)
	if err != nil {
		return err
	}
	return nil
}

func WriteUTF8BOM(dst io.Writer) error {
	_, err := dst.Write(UTF8BOM[0:3])
	if err != nil {
		return err
	}
	return nil
}

func RuneToUTF16Bytes(r rune) []byte {
	utf16Buf := utf16.Encode([]rune{r})
	b := (*[2]byte)(unsafe.Pointer(&utf16Buf[0]))
	return b[0:2]
}

func UTF8ToUTF16Ctx(ctx context.Context, src io.Reader, dst io.Writer) error {
	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(dst)

	if err := WriteUTF16BOM(byteorder.Get(), writer); err != nil {
		return err
	}

	// Read first rune and check if it is UTF-8 BOM.
	r, _, err := reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	// If first rune is NOT UTF-8 BOM(0xEF,0xBB,0xBF -> rune: 0xFEFF),
	// convert it to UTF-16 bytes, write the bytes.
	if r != 0xFEFF {
		b := RuneToUTF16Bytes(r)
		if _, err := writer.Write(b); err != nil {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		b := RuneToUTF16Bytes(r)
		if _, err := writer.Write(b); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func UTF8ToUTF16(src io.Reader, dst io.Writer) error {
	return UTF8ToUTF16Ctx(context.Background(), src, dst)
}

func UTF16ToUTF8Ctx(ctx context.Context, src io.Reader, dst io.Writer, outputUTF8BOM bool) error {
	var (
		buf       = make([]byte, 2)
		utf8Buf   = make([]byte, 4)
		uint16Buf = make([]uint16, 1)
	)

	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(dst)

	// Output UTF-8 BOM if need(E.g. Plain TXT files for Windows)
	if outputUTF8BOM {
		if err := WriteUTF8BOM(writer); err != nil {
			return err
		}
	}

	order, err := DetectUTF16BOM(reader)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		i := 0
		for i = 0; i < 2; i++ {
			buf[i], err = reader.ReadByte()
		}

		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			break
		}

		uint16Buf[0] = order.Uint16(buf[0:])
		runes := utf16.Decode(uint16Buf)

		for _, r := range runes {
			n := utf8.EncodeRune(utf8Buf, r)
			dst.Write(utf8Buf[:n])
		}
	}

	return writer.Flush()
}

func UTF16ToUTF8(src io.Reader, dst io.Writer, outputUTF8BOM bool) error {
	return UTF16ToUTF8Ctx(context.Background(), src, dst, outputUTF8BOM)
}
