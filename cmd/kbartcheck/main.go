package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miku/holdingfile/kbart"
)

func main() {
	var r io.Reader

	skipHeader := flag.Bool("skip", false, "skip header row")

	flag.Parse()

	if flag.NArg() == 0 {
		r = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		r = file
	}

	kr := kbart.NewReader(r)
	kr.SkipFirstRow = *skipHeader

	stats := make(map[string]int)
	var i int

	for {
		_, _, err := kr.Read()
		if err == io.EOF {
			break
		}
		i++
		if err != nil {
			stats[err.Error()]++
		}
	}
	stats["records"] = i
	b, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
