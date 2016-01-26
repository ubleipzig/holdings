package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miku/holdingfile/kbart"
)

func main() {
	var r io.Reader

	flag.Parse()
	if flag.NArg() == 0 {
		r = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		r = bufio.NewReader(file)
	}

	kr := kbart.NewReader(r)
	kr.SkipFirstRow = true

	for {
		_, _, err := kr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
	}

}
