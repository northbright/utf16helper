# utf16helper

[![Build Status](https://travis-ci.org/northbright/utf16helper.svg?branch=master)](https://travis-ci.org/northbright/utf16helper)
[![GoDoc](https://godoc.org/github.com/northbright/utf16helper?status.svg)](https://godoc.org/github.com/northbright/utf16helper)

utf16helper is a [Golang](https://golang.org) package provide UTF16 helper functions.
e.g. converts UTF-16 from / to UTF-8.

## Example
```
package main

import (
        "context"
        "log"
        "bytes"

        "github.com/northbright/utf16helper"
)

func main() {
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

        // Output: 48656C6C6F2C20E4B896E7958C21
        log.Printf("%X", writer.Bytes())
        // Output string converted from UTF8 bytes: Hello, 世界!
        log.Printf("string: %s", writer.String())
}
```

## Documentation
* [API References](https://godoc.org/github.com/northbright/utf16helper)

## License
* [MIT License](./LICENSE)

