package main

import (
	"flag"
	"log"

	"github.com/miku/holdingfile"
)

func main() {
	date := flag.String("date", "", "record date")
	file := flag.String("file", "", "holding file")
	issn := flag.String("issn", "", "record issn")
	issue := flag.String("issue", "", "record issue")
	volume := flag.String("volume", "", "record volume")

	flag.Parse()

	if *issn == "" {
		log.Fatal("ISSN is required")
	}

	s := holdingfile.Signature{Date: *date, Volume: *volume, Issue: *issue}
	log.Printf("test if %s %+v is covered by %s", *issn, s, *file)
}
