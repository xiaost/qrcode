// go-qrcode
// Copyright 2014 Tom Harwood
// Copyright 2025 Kyle Xiao (xiaost7@gmail.com)

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	qrcode "github.com/xiaost/qrcode"
)

func main() {
	outFile := flag.String("o", "", "out PNG file prefix, empty for stdout")
	size := flag.Int("s", 256, "image size (pixel)")
	textArt := flag.Bool("t", false, "print as text-art on stdout")
	negative := flag.Bool("i", false, "invert black and white")
	disableBorder := flag.Bool("d", false, "disable QR Code border")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `qrcode -- QR Code encoder in Go
https://github.com/xiaost/qrcode

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Usage:
  1. Arguments except for flags are joined by " " and used to generate QR code.
     Default output is STDOUT, pipe to imagemagick command "display" to display
     on any X server.

       qrcode hello word | display

  2. Save to file if "display" not available:

       qrcode "homepage: https://github.com/xiaost/qrcode" > out.png

`)
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		checkError(fmt.Errorf("Error: no content given"))
	}

	content := strings.Join(flag.Args(), " ")

	q, err := qrcode.New(content, qrcode.Highest)
	checkError(err)

	if *disableBorder {
		q.DisableBorder = true
	}

	if *textArt {
		art := q.ToString(*negative)
		fmt.Println(art)
		return
	}

	if *negative {
		q.ForegroundColor, q.BackgroundColor = q.BackgroundColor, q.ForegroundColor
	}

	png, err := q.PNG(*size)
	checkError(err)

	if *outFile == "" {
		_, err = os.Stdout.Write(png)
		checkError(err)
	} else {
		fh, err := os.Create(*outFile + ".png")
		checkError(err)
		defer fh.Close()
		_, err = fh.Write(png)
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
